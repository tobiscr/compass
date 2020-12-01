package eventdef

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	log "github.com/sirupsen/logrus"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/pkg/errors"
)

const eventAPIDefTable string = `"public"."event_api_definitions"`

var (
	idColumn      = "id"
	tenantColumn  = "tenant_id"
	packageColumn = "package_id"
	apiDefColumns = []string{idColumn, tenantColumn, packageColumn, "name", "description", "group_name", "spec_data",
		"spec_format", "spec_type", "version_value", "version_deprecated", "version_deprecated_since",
		"version_for_removal"}
	idColumns        = []string{"id"}
	updatableColumns = []string{"name", "description", "group_name", "spec_data", "spec_format", "spec_type",
		"version_value", "version_deprecated", "version_deprecated_since", "version_for_removal"}
)

//go:generate mockery -name=EventAPIDefinitionConverter -output=automock -outpkg=automock -case=underscore
type EventAPIDefinitionConverter interface {
	FromEntity(entity Entity) (model.EventDefinition, error)
	ToEntity(apiModel model.EventDefinition) (Entity, error)
}

type pgRepository struct {
	singleGetter    repo.SingleGetter
	singleLister    repo.Lister
	pageableQuerier repo.PageableQuerier
	creator         repo.Creator
	updater         repo.Updater
	deleter         repo.Deleter
	existQuerier    repo.ExistQuerier
	conv            EventAPIDefinitionConverter
}

func NewRepository(conv EventAPIDefinitionConverter) *pgRepository {
	return &pgRepository{
		singleGetter:    repo.NewSingleGetter(resource.EventDefinition, eventAPIDefTable, tenantColumn, apiDefColumns),
		singleLister:    repo.NewLister(resource.EventDefinition, eventAPIDefTable, tenantColumn, apiDefColumns),
		pageableQuerier: repo.NewPageableQuerier(resource.EventDefinition, eventAPIDefTable, tenantColumn, apiDefColumns),
		creator:         repo.NewCreator(resource.EventDefinition, eventAPIDefTable, apiDefColumns),
		updater:         repo.NewUpdater(resource.EventDefinition, eventAPIDefTable, updatableColumns, tenantColumn, idColumns),
		deleter:         repo.NewDeleter(resource.EventDefinition, eventAPIDefTable, tenantColumn),
		existQuerier:    repo.NewExistQuerier(resource.EventDefinition, eventAPIDefTable, tenantColumn),
		conv:            conv,
	}
}

type EventAPIDefCollection []Entity

func (r EventAPIDefCollection) Len() int {
	return len(r)
}

func (r *pgRepository) GetByID(ctx context.Context, tenantID string, id string) (*model.EventDefinition, error) {
	var eventAPIDefEntity Entity
	err := r.singleGetter.Get(ctx, tenantID, repo.Conditions{repo.NewEqualCondition("id", id)}, repo.NoOrderBy, &eventAPIDefEntity)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting EventDefinition with id %s", id)
	}

	eventAPIDefModel, err := r.conv.FromEntity(eventAPIDefEntity)
	if err != nil {
		return nil, errors.Wrap(err, "while creating EventDefinition entity to model")
	}

	return &eventAPIDefModel, nil
}

func (r *pgRepository) GetForPackage(ctx context.Context, tenant string, id string, packageID string) (*model.EventDefinition, error) {
	var ent Entity

	conditions := repo.Conditions{
		repo.NewEqualCondition(idColumn, id),
		repo.NewEqualCondition(packageColumn, packageID),
	}
	if err := r.singleGetter.Get(ctx, tenant, conditions, repo.NoOrderBy, &ent); err != nil {
		return nil, err
	}

	eventAPIModel, err := r.conv.FromEntity(ent)
	if err != nil {
		return nil, errors.Wrap(err, "while creating event definition model from entity")
	}

	return &eventAPIModel, nil
}

func (r *pgRepository) ListForPackage(ctx context.Context, tenantID string, packageID string, pageSize int, cursor string) (*model.EventDefinitionPage, error) {
	conditions := repo.Conditions{
		repo.NewEqualCondition(packageColumn, packageID),
	}

	return r.list(ctx, tenantID, pageSize, cursor, conditions)
}

func (r *pgRepository) ListAllForPackage(ctx context.Context, tenantID string, packageIDs []string, pageSize int, cursor string) ([]*model.EventDefinitionPage, error) {
	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var eventDefCollection EventAPIDefCollection
	var query string
	var sb strings.Builder
	unionQuery := fmt.Sprint("(SELECT id, tenant_id, package_id, name, description, group_name, spec_data, spec_format, spec_type, version_value, version_deprecated, version_deprecated_since, version_for_removal from event_api_definitions WHERE package_id='%s' and tenant_id='%s' order by id limit %v offset %d)")

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
	err = persist.Select(&eventDefCollection, query)
	if err != nil {
		return nil, err
	}

	// count logic
	conditions := repo.Conditions{
		repo.NewInConditionForStringValues("package_id", packageIDs),
	}
	var eventDefCollectionCount EventAPIDefCollection
	err = r.singleLister.List(ctx, tenantID, &eventDefCollectionCount, conditions...)
	if err != nil {
		return nil, err
	}

	eventDefsCountByID := map[string][]*model.EventDefinition{}
	for _, eventDefEnt := range eventDefCollectionCount {
		m, err := r.conv.FromEntity(eventDefEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating EventDef model from entity")
		}
		eventDefsCountByID[eventDefEnt.PkgID] = append(eventDefsCountByID[eventDefEnt.PkgID], &m)
	}
	// end

	eventDefsById := map[string][]*model.EventDefinition{}
	for _, eventDefEnt := range eventDefCollection {
		m, err := r.conv.FromEntity(eventDefEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating EventDef model from entity")
		}
		eventDefsById[eventDefEnt.PkgID] = append(eventDefsById[eventDefEnt.PkgID], &m)
	}

	// map the ApiDefPage to the current package_id
	eventDefPages := make([]*model.EventDefinitionPage, len(packageIDs))
	for i, pkgID := range packageIDs {
		totalCount := len(eventDefsCountByID[pkgID])
		hasNextPage := false
		endCursor := ""
		if totalCount > offset+len(eventDefsById[pkgID]) {
			hasNextPage = true
			endCursor = pagination.EncodeNextOffsetCursor(offset, pageSize)
		}

		page := &pagination.Page{
			StartCursor: cursor,
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		}

		eventDefPages[i] = &model.EventDefinitionPage{Data: eventDefsById[pkgID], TotalCount: totalCount, PageInfo: page}
	}

	return eventDefPages, nil
}

func (r *pgRepository) ListAllForPackageNoPaging(ctx context.Context, tenantID string, packageIDs []string) ([][]*model.EventDefinition, error) {
	conditions := repo.Conditions{
		repo.NewInConditionForStringValues("package_id", packageIDs),
	}

	var eventDefCollection EventAPIDefCollection
	err := r.singleLister.List(ctx, tenantID, &eventDefCollection, conditions...)
	if err != nil {
		return nil, err
	}

	eventDefsByID := map[string][]*model.EventDefinition{}
	for _, eventDefEnt := range eventDefCollection {
		m, err := r.conv.FromEntity(eventDefEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating EventDef model from entity")
		}
		eventDefsByID[eventDefEnt.PkgID] = append(eventDefsByID[eventDefEnt.PkgID], &m)
	}

	// map the PackagePage to the current pkg_id
	eventDefs := make([][]*model.EventDefinition, len(packageIDs))
	for i, pkgID := range packageIDs {
		eventDefs[i] = eventDefsByID[pkgID]
	}

	return eventDefs, nil
}

func (r *pgRepository) list(ctx context.Context, tenant string, pageSize int, cursor string, conditions repo.Conditions) (*model.EventDefinitionPage, error) {
	var eventCollection EventAPIDefCollection
	page, totalCount, err := r.pageableQuerier.List(ctx, tenant, pageSize, cursor, idColumn, &eventCollection, conditions...)
	if err != nil {
		return nil, err
	}

	var items []*model.EventDefinition

	for _, eventEnt := range eventCollection {
		m, err := r.conv.FromEntity(eventEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating APIDefinition model from entity")
		}
		items = append(items, &m)
	}

	return &model.EventDefinitionPage{
		Data:       items,
		TotalCount: totalCount,
		PageInfo:   page,
	}, nil
}

func (r *pgRepository) Create(ctx context.Context, item *model.EventDefinition) error {
	if item == nil {
		return apperrors.NewInternalError("item cannot be nil")
	}

	entity, err := r.conv.ToEntity(*item)
	if err != nil {
		return errors.Wrap(err, "while creating EventDefinition model to entity")
	}

	log.Debugf("Persisting Event-Definition entity with id %s to db", item.ID)
	err = r.creator.Create(ctx, entity)
	if err != nil {
		return errors.Wrap(err, "while saving entity to db")
	}

	return nil
}

func (r *pgRepository) CreateMany(ctx context.Context, items []*model.EventDefinition) error {
	for index, item := range items {
		entity, err := r.conv.ToEntity(*item)
		if err != nil {
			return errors.Wrapf(err, "while creating %d item", index)
		}
		err = r.creator.Create(ctx, entity)
		if err != nil {
			return errors.Wrapf(err, "while persisting %d item", index)
		}
	}

	return nil
}

func (r *pgRepository) Update(ctx context.Context, item *model.EventDefinition) error {
	if item == nil {
		return apperrors.NewInternalError("item cannot be nil")
	}

	entity, err := r.conv.ToEntity(*item)
	if err != nil {
		return errors.Wrap(err, "while converting model to entity")
	}

	return r.updater.UpdateSingle(ctx, entity)
}

func (r *pgRepository) Exists(ctx context.Context, tenantID, id string) (bool, error) {
	return r.existQuerier.Exists(ctx, tenantID, repo.Conditions{repo.NewEqualCondition(idColumn, id)})
}

func (r *pgRepository) Delete(ctx context.Context, tenantID string, id string) error {
	return r.deleter.DeleteOne(ctx, tenantID, repo.Conditions{repo.NewEqualCondition(idColumn, id)})
}
