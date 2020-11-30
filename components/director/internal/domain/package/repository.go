package mp_package

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"

	log "github.com/sirupsen/logrus"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/pkg/errors"
)

const packageTable string = `public.packages`
const packageInstanceAuthTable string = `public.package_instance_auths`
const packageInstanceAuthPackageRefField string = `package_id`

var (
	packageColumns = []string{"id", "tenant_id", "app_id", "name", "description", "instance_auth_request_json_schema", "default_instance_auth"}
	tenantColumn   = "tenant_id"
)

//go:generate mockery -name=EntityConverter -output=automock -outpkg=automock -case=underscore
type EntityConverter interface {
	ToEntity(in *model.Package) (*Entity, error)
	FromEntity(entity *Entity) (*model.Package, error)
}

type pgRepository struct {
	existQuerier    repo.ExistQuerier
	singleGetter    repo.SingleGetter
	singleLister    repo.Lister
	deleter         repo.Deleter
	pageableQuerier repo.PageableQuerier
	creator         repo.Creator
	updater         repo.Updater
	conv            EntityConverter
}

func NewRepository(conv EntityConverter) *pgRepository {
	return &pgRepository{
		existQuerier:    repo.NewExistQuerier(resource.Package, packageTable, tenantColumn),
		singleGetter:    repo.NewSingleGetter(resource.Package, packageTable, tenantColumn, packageColumns),
		singleLister:    repo.NewLister(resource.Package, packageTable, tenantColumn, packageColumns),
		deleter:         repo.NewDeleter(resource.Package, packageTable, tenantColumn),
		pageableQuerier: repo.NewPageableQuerier(resource.Package, packageTable, tenantColumn, packageColumns),
		creator:         repo.NewCreator(resource.Package, packageTable, packageColumns),
		updater:         repo.NewUpdater(resource.Package, packageTable, []string{"name", "description", "instance_auth_request_json_schema", "default_instance_auth"}, tenantColumn, []string{"id"}),
		conv:            conv,
	}
}

type PackageCollection []Entity

func (r PackageCollection) Len() int {
	return len(r)
}

func (r *pgRepository) Create(ctx context.Context, model *model.Package) error {
	if model == nil {
		return apperrors.NewInternalError("model can not be nil")
	}

	pkgEnt, err := r.conv.ToEntity(model)
	if err != nil {
		return errors.Wrap(err, "while converting to Package entity")
	}

	log.Debugf("Persisting Package entity with id %s to db", model.ID)
	return r.creator.Create(ctx, pkgEnt)
}

func (r *pgRepository) Update(ctx context.Context, model *model.Package) error {
	if model == nil {
		return apperrors.NewInternalError("model can not be nil")
	}

	pkgEnt, err := r.conv.ToEntity(model)

	if err != nil {
		return errors.Wrap(err, "while converting to Package entity")
	}

	return r.updater.UpdateSingle(ctx, pkgEnt)
}

func (r *pgRepository) Delete(ctx context.Context, tenant, id string) error {
	return r.deleter.DeleteOne(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)})
}

func (r *pgRepository) Exists(ctx context.Context, tenant, id string) (bool, error) {
	return r.existQuerier.Exists(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)})
}

func (r *pgRepository) GetByID(ctx context.Context, tenant, id string) (*model.Package, error) {
	var pkgEnt Entity
	if err := r.singleGetter.Get(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)}, repo.NoOrderBy, &pkgEnt); err != nil {
		return nil, err
	}

	pkgModel, err := r.conv.FromEntity(&pkgEnt)
	if err != nil {
		return nil, errors.Wrap(err, "while converting Package from Entity")
	}

	return pkgModel, nil
}

func (r *pgRepository) GetForApplication(ctx context.Context, tenant string, id string, applicationID string) (*model.Package, error) {
	var ent Entity

	conditions := repo.Conditions{
		repo.NewEqualCondition("id", id),
		repo.NewEqualCondition("app_id", applicationID),
	}
	if err := r.singleGetter.Get(ctx, tenant, conditions, repo.NoOrderBy, &ent); err != nil {
		return nil, err
	}

	pkgModel, err := r.conv.FromEntity(&ent)
	if err != nil {
		return nil, errors.Wrap(err, "while creating Package model from entity")
	}

	return pkgModel, nil
}

func (r *pgRepository) GetByInstanceAuthID(ctx context.Context, tenant string, instanceAuthID string) (*model.Package, error) {
	var pkgEnt Entity

	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return nil, err
	}

	prefixedFieldNames := str.PrefixStrings(packageColumns, "p.")
	stmt := fmt.Sprintf(`SELECT %s FROM %s AS p JOIN %s AS a on a.%s=p.id where a.tenant_id=$1 AND a.id=$2`,
		strings.Join(prefixedFieldNames, ", "),
		packageTable,
		packageInstanceAuthTable,
		packageInstanceAuthPackageRefField)

	err = persist.Get(&pkgEnt, stmt, tenant, instanceAuthID)
	switch {
	case err != nil:
		return nil, errors.Wrap(err, "while getting Package by Instance Auth ID")
	}

	pkgModel, err := r.conv.FromEntity(&pkgEnt)
	if err != nil {
		return nil, errors.Wrap(err, "while creating Package model from entity")
	}

	return pkgModel, nil
}

func (r *pgRepository) ListByApplicationID(ctx context.Context, tenantID string, applicationID string, pageSize int, cursor string) (*model.PackagePage, error) {
	conditions := repo.Conditions{
		repo.NewEqualCondition("app_id", applicationID),
	}

	var packageCollection PackageCollection
	page, totalCount, err := r.pageableQuerier.List(ctx, tenantID, pageSize, cursor, "id", &packageCollection, conditions...)
	if err != nil {
		return nil, err
	}

	var items []*model.Package

	for _, pkgEnt := range packageCollection {
		m, err := r.conv.FromEntity(&pkgEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Package model from entity")
		}
		items = append(items, m)
	}

	return &model.PackagePage{
		Data:       items,
		TotalCount: totalCount,
		PageInfo:   page,
	}, nil
}

func (r *pgRepository) ListByApplicationIDs(ctx context.Context, tenantID string, applicationIDs []string, pageSize int, cursor string) ([]*model.PackagePage, error) {
	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var packageCollection PackageCollection
	var query string
	var sb strings.Builder
	unionQuery := fmt.Sprint("(SELECT name, id, app_id from packages WHERE app_id='%s' and tenant_id='%s' order by id limit %v offset %d)")

	offset, err := pagination.DecodeOffsetCursor(cursor)
	if err != nil {
		return nil, errors.Wrap(err, "while decoding page cursor")
	}

	for i := 0; i < len(applicationIDs); i++ {
		if i == len(applicationIDs)-1 {
			sb.WriteString(fmt.Sprintf(unionQuery, applicationIDs[i], tenantID, pageSize, offset))
			query = sb.String()
		}
		sb.WriteString(fmt.Sprintf(unionQuery, applicationIDs[i], tenantID, pageSize, offset) + "union")
	}

	fmt.Println("[===Executing single union query ===] ", query)
	err = persist.Select(&packageCollection, query)
	if err != nil {
		return nil, err
	}

	//count logic
	conditions := repo.Conditions{
		repo.NewInConditionForStringValues("app_id", applicationIDs),
	}
	var packageCollectionCount PackageCollection
	err = r.singleLister.List(ctx, tenantID, &packageCollectionCount, conditions...)
	if err != nil {
		return nil, err
	}

	pkgsCountByID := map[string][]*model.Package{}
	for _, pkgEnt := range packageCollectionCount {
		m, err := r.conv.FromEntity(&pkgEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Package model from entity")
		}
		pkgsCountByID[pkgEnt.ApplicationID] = append(pkgsCountByID[pkgEnt.ApplicationID], m)
	}
	//end

	pkgByID := map[string][]*model.Package{}
	for _, pkgEnt := range packageCollection {
		m, err := r.conv.FromEntity(&pkgEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Package model from entity")
		}
		pkgByID[pkgEnt.ApplicationID] = append(pkgByID[pkgEnt.ApplicationID], m)
	}

	// map the PackagePage to the current app_id
	pkgPages := make([]*model.PackagePage, len(applicationIDs))
	for i, appID := range applicationIDs {
		totalCount := len(pkgsCountByID[appID])
		hasNextPage := false
		endCursor := ""
		if totalCount > offset+len(pkgByID[appID]) {
			hasNextPage = true
			endCursor = pagination.EncodeNextOffsetCursor(offset, pageSize)
		}

		page := &pagination.Page{
			StartCursor: cursor,
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		}

		pkgPages[i] = &model.PackagePage{Data: pkgByID[appID], TotalCount: totalCount, PageInfo: page}
	}

	return pkgPages, nil
}

func (r *pgRepository) ListByApplicationIDsNoPaging(ctx context.Context, tnt string, applicationIDs []string) ([][]*model.Package, error) {
	conditions := repo.Conditions{
		repo.NewInConditionForStringValues("app_id", applicationIDs),
	}

	var packageCollection PackageCollection
	err := r.singleLister.List(ctx, tnt, &packageCollection, conditions...)
	if err != nil {
		return nil, err
	}

	pkgsByID := map[string][]*model.Package{}
	for _, pkgEnt := range packageCollection {
		m, err := r.conv.FromEntity(&pkgEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Package model from entity")
		}
		pkgsByID[pkgEnt.ApplicationID] = append(pkgsByID[pkgEnt.ApplicationID], m)
	}

	// map the PackagePage to the current app_id
	pkgs := make([][]*model.Package, len(applicationIDs))
	for i, appID := range applicationIDs {
		pkgs[i] = pkgsByID[appID]
	}

	return pkgs, nil
}
