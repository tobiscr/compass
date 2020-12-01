package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/pkg/errors"
)

const apiDefTable string = `"public"."api_definitions"`

var (
	tenantColumn  = "tenant_id"
	apiDefColumns = []string{"id", "tenant_id", "package_id", "name", "description", "group_name", "target_url", "spec_data",
		"spec_format", "spec_type", "version_value", "version_deprecated", "version_deprecated_since", "version_for_removal"}
	idColumns        = []string{"id"}
	updatableColumns = []string{"name", "description", "group_name", "target_url", "spec_data", "spec_format", "spec_type",
		"version_value", "version_deprecated", "version_deprecated_since", "version_for_removal"}
)

//go:generate mockery -name=APIDefinitionConverter -output=automock -outpkg=automock -case=underscore
type APIDefinitionConverter interface {
	FromEntity(entity Entity) model.APIDefinition
	ToEntity(apiModel model.APIDefinition) Entity
}

type pgRepository struct {
	creator         repo.Creator
	singleGetter    repo.SingleGetter
	singleLister    repo.Lister
	pageableQuerier repo.PageableQuerier
	updater         repo.Updater
	deleter         repo.Deleter
	existQuerier    repo.ExistQuerier
	conv            APIDefinitionConverter
}

func NewRepository(conv APIDefinitionConverter) *pgRepository {
	return &pgRepository{
		singleGetter:    repo.NewSingleGetter(resource.API, apiDefTable, tenantColumn, apiDefColumns),
		singleLister:    repo.NewLister(resource.API, apiDefTable, tenantColumn, apiDefColumns),
		pageableQuerier: repo.NewPageableQuerier(resource.API, apiDefTable, tenantColumn, apiDefColumns),
		creator:         repo.NewCreator(resource.API, apiDefTable, apiDefColumns),
		updater:         repo.NewUpdater(resource.API, apiDefTable, updatableColumns, tenantColumn, idColumns),
		deleter:         repo.NewDeleter(resource.API, apiDefTable, tenantColumn),
		existQuerier:    repo.NewExistQuerier(resource.API, apiDefTable, tenantColumn),
		conv:            conv,
	}
}

type APIDefCollection []Entity

func (r APIDefCollection) Len() int {
	return len(r)
}

func (r *pgRepository) ListForPackage(ctx context.Context, tenantID string, packageID string, pageSize int, cursor string) (*model.APIDefinitionPage, error) {
	conditions := repo.Conditions{
		repo.NewEqualCondition("package_id", packageID),
	}
	return r.list(ctx, tenantID, pageSize, cursor, conditions)
}

func (r *pgRepository) ListAllForPackage(ctx context.Context, tenantID string, packageIDs []string, pageSize int, cursor string) ([]*model.APIDefinitionPage, error) {
	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var apiDefCollection APIDefCollection
	var query string
	var sb strings.Builder
	unionQuery := fmt.Sprint("(SELECT id, tenant_id, package_id, name, description, group_name, target_url, spec_data, spec_format, spec_type, version_value, version_deprecated, version_deprecated_since, version_for_removal from api_definitions WHERE package_id='%s' and tenant_id='%s' order by id limit %v offset %d)")

	offset, err := pagination.DecodeOffsetCursor(cursor)
	if err != nil {
		return nil, errors.Wrap(err, "while decoding page cursor")
	}

	for i := 0; i < len(packageIDs); i++ {
		if i == len(packageIDs)-1 {
			sb.WriteString(fmt.Sprintf(unionQuery, packageIDs[i], tenantID, pageSize, offset))
			query = sb.String()
		}
		sb.WriteString(fmt.Sprintf(unionQuery, packageIDs[i], tenantID, pageSize, offset) + "union")
	}

	fmt.Println("[===Executing single union query ===] ", query)
	err = persist.Select(&apiDefCollection, query)
	if err != nil {
		return nil, err
	}

	//count logic
	conditions := repo.Conditions{
		repo.NewInConditionForStringValues("package_id", packageIDs),
	}
	var apiDefCollectionCount APIDefCollection
	err = r.singleLister.List(ctx, tenantID, &apiDefCollectionCount, conditions...)
	if err != nil {
		return nil, err
	}

	apiDefsCountByID := map[string][]*model.APIDefinition{}
	for _, apiDefEnt := range apiDefCollectionCount {
		m := r.conv.FromEntity(apiDefEnt)
		apiDefsCountByID[apiDefEnt.PkgID] = append(apiDefsCountByID[apiDefEnt.PkgID], &m)
	}
	//end

	apiDefsById := map[string][]*model.APIDefinition{}
	for _, apiDefEnt := range apiDefCollection {
		m := r.conv.FromEntity(apiDefEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating ApiDef model from entity")
		}
		apiDefsById[apiDefEnt.PkgID] = append(apiDefsById[apiDefEnt.PkgID], &m)
	}

	// map the ApiDefPage to the current package_id
	apiDefPages := make([]*model.APIDefinitionPage, len(packageIDs))
	for i, pkgID := range packageIDs {
		totalCount := len(apiDefsCountByID[pkgID])
		hasNextPage := false
		endCursor := ""
		if totalCount > offset+len(apiDefsById[pkgID]) {
			hasNextPage = true
			endCursor = pagination.EncodeNextOffsetCursor(offset, pageSize)
		}

		page := &pagination.Page{
			StartCursor: cursor,
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		}

		apiDefPages[i] = &model.APIDefinitionPage{Data: apiDefsById[pkgID], TotalCount: totalCount, PageInfo: page}
	}

	return apiDefPages, nil
}

func (r *pgRepository) ListAllForPackageNoPaging(ctx context.Context, tenantID string, packageIDs []string) ([][]*model.APIDefinition, error) {
	conditions := repo.Conditions{
		repo.NewInConditionForStringValues("package_id", packageIDs),
	}

	var apiDefCollection APIDefCollection
	err := r.singleLister.List(ctx, tenantID, &apiDefCollection, conditions...)
	if err != nil {
		return nil, err
	}

	apiDefsByID := map[string][]*model.APIDefinition{}
	for _, apiDefEnt := range apiDefCollection {
		m := r.conv.FromEntity(apiDefEnt)
		apiDefsByID[apiDefEnt.PkgID] = append(apiDefsByID[apiDefEnt.PkgID], &m)
	}

	// map the PackagePage to the current pkg_id
	apiDefs := make([][]*model.APIDefinition, len(packageIDs))
	for i, pkgID := range packageIDs {
		apiDefs[i] = apiDefsByID[pkgID]
	}

	return apiDefs, nil
}

func (r *pgRepository) list(ctx context.Context, tenant string, pageSize int, cursor string, conditions repo.Conditions) (*model.APIDefinitionPage, error) {
	var apiDefCollection APIDefCollection
	page, totalCount, err := r.pageableQuerier.List(ctx, tenant, pageSize, cursor, "id", &apiDefCollection, conditions...)
	if err != nil {
		return nil, err
	}

	var items []*model.APIDefinition

	for _, apiDefEnt := range apiDefCollection {
		m := r.conv.FromEntity(apiDefEnt)

		items = append(items, &m)
	}

	return &model.APIDefinitionPage{
		Data:       items,
		TotalCount: totalCount,
		PageInfo:   page,
	}, nil
}

func (r *pgRepository) GetByID(ctx context.Context, tenantID string, id string) (*model.APIDefinition, error) {
	var apiDefEntity Entity
	err := r.singleGetter.Get(ctx, tenantID, repo.Conditions{repo.NewEqualCondition("id", id)}, repo.NoOrderBy, &apiDefEntity)
	if err != nil {
		return nil, errors.Wrap(err, "while getting APIDefinition")
	}

	apiDefModel := r.conv.FromEntity(apiDefEntity)

	return &apiDefModel, nil
}

func (r *pgRepository) GetForPackage(ctx context.Context, tenant string, id string, packageID string) (*model.APIDefinition, error) {
	var ent Entity

	conditions := repo.Conditions{
		repo.NewEqualCondition("id", id),
		repo.NewEqualCondition("package_id", packageID),
	}
	if err := r.singleGetter.Get(ctx, tenant, conditions, repo.NoOrderBy, &ent); err != nil {
		return nil, err
	}

	apiDefModel := r.conv.FromEntity(ent)

	return &apiDefModel, nil
}

func (r *pgRepository) Create(ctx context.Context, item *model.APIDefinition) error {
	if item == nil {
		return apperrors.NewInternalError("item cannot be nil")
	}

	entity := r.conv.ToEntity(*item)
	err := r.creator.Create(ctx, entity)
	if err != nil {
		return errors.Wrap(err, "while saving entity to db")
	}

	return nil
}

func (r *pgRepository) CreateMany(ctx context.Context, items []*model.APIDefinition) error {
	for index, item := range items {
		entity := r.conv.ToEntity(*item)

		err := r.creator.Create(ctx, entity)
		if err != nil {
			return errors.Wrapf(err, "while persisting %d item", index)
		}
	}

	return nil
}

func (r *pgRepository) Update(ctx context.Context, item *model.APIDefinition) error {
	if item == nil {
		return apperrors.NewInternalError("item cannot be nil")
	}

	entity := r.conv.ToEntity(*item)

	return r.updater.UpdateSingle(ctx, entity)
}

func (r *pgRepository) Exists(ctx context.Context, tenantID, id string) (bool, error) {
	return r.existQuerier.Exists(ctx, tenantID, repo.Conditions{repo.NewEqualCondition("id", id)})
}

func (r *pgRepository) Delete(ctx context.Context, tenantID string, id string) error {
	return r.deleter.DeleteOne(ctx, tenantID, repo.Conditions{repo.NewEqualCondition("id", id)})
}
