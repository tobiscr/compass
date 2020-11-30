package mp_package

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/dataloader"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/internal/model"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

//go:generate mockery -name=PackageService -output=automock -outpkg=automock -case=underscore
type PackageService interface {
	Create(ctx context.Context, applicationID string, in model.PackageCreateInput) (string, error)
	Update(ctx context.Context, id string, in model.PackageUpdateInput) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*model.Package, error)
}

//go:generate mockery -name=PackageConverter -output=automock -outpkg=automock -case=underscore
type PackageConverter interface {
	ToGraphQL(in *model.Package) (*graphql.Package, error)
	CreateInputFromGraphQL(in graphql.PackageCreateInput) (model.PackageCreateInput, error)
	UpdateInputFromGraphQL(in graphql.PackageUpdateInput) (*model.PackageUpdateInput, error)
}

//go:generate mockery -name=PackageInstanceAuthService -output=automock -outpkg=automock -case=underscore
type PackageInstanceAuthService interface {
	GetForPackage(ctx context.Context, id string, packageID string) (*model.PackageInstanceAuth, error)
	List(ctx context.Context, id string) ([]*model.PackageInstanceAuth, error)
}

//go:generate mockery -name=PackageInstanceAuthConverter -output=automock -outpkg=automock -case=underscore
type PackageInstanceAuthConverter interface {
	ToGraphQL(in *model.PackageInstanceAuth) (*graphql.PackageInstanceAuth, error)
	MultipleToGraphQL(in []*model.PackageInstanceAuth) ([]*graphql.PackageInstanceAuth, error)
}

//go:generate mockery -name=APIService -output=automock -outpkg=automock -case=underscore
type APIService interface {
	ListForPackage(ctx context.Context, packageID string, pageSize int, cursor string) (*model.APIDefinitionPage, error)
	ListAllByPackageIDs(ctx context.Context, packageIDs []string, pageSize int, cursor string) ([]*model.APIDefinitionPage, error)
	ListAllByPackageIDsNoPaging(ctx context.Context, packageIDs []string) ([][]*model.APIDefinition, error)
	GetForPackage(ctx context.Context, id string, packageID string) (*model.APIDefinition, error)
}

//go:generate mockery -name=APIConverter -output=automock -outpkg=automock -case=underscore
type APIConverter interface {
	ToGraphQL(in *model.APIDefinition) *graphql.APIDefinition
	MultipleToGraphQL(in []*model.APIDefinition) []*graphql.APIDefinition
	MultipleInputFromGraphQL(in []*graphql.APIDefinitionInput) ([]*model.APIDefinitionInput, error)
}

//go:generate mockery -name=EventService -output=automock -outpkg=automock -case=underscore
type EventService interface {
	ListForPackage(ctx context.Context, packageID string, pageSize int, cursor string) (*model.EventDefinitionPage, error)
	ListAllByPackageIDs(ctx context.Context, packageIDs []string, pageSize int, cursor string) ([]*model.EventDefinitionPage, error)
	ListAllByPackageIDsNoPaging(ctx context.Context, packageIDs []string) ([][]*model.EventDefinition, error)
	GetForPackage(ctx context.Context, id string, packageID string) (*model.EventDefinition, error)
}

//go:generate mockery -name=EventConverter -output=automock -outpkg=automock -case=underscore
type EventConverter interface {
	ToGraphQL(in *model.EventDefinition) *graphql.EventDefinition
	MultipleToGraphQL(in []*model.EventDefinition) []*graphql.EventDefinition
	MultipleInputFromGraphQL(in []*graphql.EventDefinitionInput) ([]*model.EventDefinitionInput, error)
}

//go:generate mockery -name=DocumentService -output=automock -outpkg=automock -case=underscore
type DocumentService interface {
	ListForPackage(ctx context.Context, packageID string, pageSize int, cursor string) (*model.DocumentPage, error)
	GetForPackage(ctx context.Context, id string, packageID string) (*model.Document, error)
}

//go:generate mockery -name=DocumentConverter -output=automock -outpkg=automock -case=underscore
type DocumentConverter interface {
	ToGraphQL(in *model.Document) *graphql.Document
	MultipleToGraphQL(in []*model.Document) []*graphql.Document
	MultipleInputFromGraphQL(in []*graphql.DocumentInput) ([]*model.DocumentInput, error)
}

type Resolver struct {
	transact persistence.Transactioner

	packageSvc             PackageService
	packageInstanceAuthSvc PackageInstanceAuthService
	apiSvc                 APIService
	eventSvc               EventService
	documentSvc            DocumentService

	packageConverter             PackageConverter
	packageInstanceAuthConverter PackageInstanceAuthConverter
	apiConverter                 APIConverter
	eventConverter               EventConverter
	documentConverter            DocumentConverter
}

func NewResolver(
	transact persistence.Transactioner,
	packageSvc PackageService,
	packageInstanceAuthSvc PackageInstanceAuthService,
	apiSvc APIService,
	eventSvc EventService,
	documentSvc DocumentService,
	packageConverter PackageConverter,
	packageInstanceAuthConverter PackageInstanceAuthConverter,
	apiConv APIConverter,
	eventConv EventConverter,
	documentConv DocumentConverter) *Resolver {
	return &Resolver{
		transact:                     transact,
		packageConverter:             packageConverter,
		packageSvc:                   packageSvc,
		packageInstanceAuthSvc:       packageInstanceAuthSvc,
		apiSvc:                       apiSvc,
		eventSvc:                     eventSvc,
		documentSvc:                  documentSvc,
		packageInstanceAuthConverter: packageInstanceAuthConverter,
		apiConverter:                 apiConv,
		eventConverter:               eventConv,
		documentConverter:            documentConv,
	}
}

func (r *Resolver) AddPackage(ctx context.Context, applicationID string, in graphql.PackageCreateInput) (*graphql.Package, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	log.Infof("Adding package to Application with id %s", applicationID)

	ctx = persistence.SaveToContext(ctx, tx)

	convertedIn, err := r.packageConverter.CreateInputFromGraphQL(in)
	if err != nil {
		return nil, errors.Wrap(err, "while converting input from GraphQL")
	}

	id, err := r.packageSvc.Create(ctx, applicationID, convertedIn)
	if err != nil {
		return nil, err
	}

	pkg, err := r.packageSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlPackage, err := r.packageConverter.ToGraphQL(pkg)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting Package with id %s to GraphQL", id)
	}

	log.Infof("Package with id %s successfully added to Application with id %s", id, applicationID)
	return gqlPackage, nil
}

func (r *Resolver) UpdatePackage(ctx context.Context, id string, in graphql.PackageUpdateInput) (*graphql.Package, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	log.Infof("Updating Package with id %s", id)

	ctx = persistence.SaveToContext(ctx, tx)

	convertedIn, err := r.packageConverter.UpdateInputFromGraphQL(in)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting converting GraphQL input to Package with id %s", id)
	}

	err = r.packageSvc.Update(ctx, id, *convertedIn)
	if err != nil {
		return nil, err
	}

	pkg, err := r.packageSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlPkg, err := r.packageConverter.ToGraphQL(pkg)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting Package with id %s to GraphQL", id)
	}

	log.Infof("Package with id %s successfully updated.", id)
	return gqlPkg, nil
}

func (r *Resolver) DeletePackage(ctx context.Context, id string) (*graphql.Package, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	log.Infof("Deleting Package with id %s", id)

	ctx = persistence.SaveToContext(ctx, tx)

	pkg, err := r.packageSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = r.packageSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	deletedPkg, err := r.packageConverter.ToGraphQL(pkg)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting Package with id %s to GraphQL", id)
	}

	log.Infof("Package with id %s successfully deleted.", id)
	return deletedPkg, nil
}

func (r *Resolver) InstanceAuth(ctx context.Context, obj *graphql.Package, id string) (*graphql.PackageInstanceAuth, error) {
	if obj == nil {
		return nil, apperrors.NewInternalError("Package cannot be empty")
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	pkg, err := r.packageInstanceAuthSvc.GetForPackage(ctx, id, obj.ID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return nil, tx.Commit()
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.packageInstanceAuthConverter.ToGraphQL(pkg)

}

func (r *Resolver) InstanceAuths(ctx context.Context, obj *graphql.Package) ([]*graphql.PackageInstanceAuth, error) {
	if obj == nil {
		return nil, apperrors.NewInternalError("Package cannot be empty")
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)
	ctx = persistence.SaveToContext(ctx, tx)

	pkgInstanceAuths, err := r.packageInstanceAuthSvc.List(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.packageInstanceAuthConverter.MultipleToGraphQL(pkgInstanceAuths)
}

func (r *Resolver) APIDefinition(ctx context.Context, obj *graphql.Package, id string) (*graphql.APIDefinition, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	api, err := r.apiSvc.GetForPackage(ctx, id, obj.ID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return nil, tx.Commit()
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.apiConverter.ToGraphQL(api), nil
}

func (r *Resolver) APIDefinitions(ctx context.Context, obj *graphql.Package, group *string, first *int, after *graphql.PageCursor) (*graphql.APIDefinitionPage, error) {
	param := dataloader.ParamApiDef{ID: obj.ID, Ctx: ctx, First: first, After: after}
	return dataloader.ApiDefFor(ctx).ApiDefById.Load(param)
}

func (r *Resolver) ApiDefinitionsDataLoader(keys []dataloader.ParamApiDef) ([]*graphql.APIDefinitionPage, []error) {
	if len(keys) == 0 {
		return nil, []error{apperrors.NewInternalError("No Packages found")}
	}

	var ctx context.Context
	var first *int
	var after *graphql.PageCursor

	packageIDs := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		if i == 0 {
			ctx = keys[i].Ctx
			first = keys[i].First
			after = keys[i].After
			packageIDs[i] = keys[i].ID
		}
		packageIDs[i] = keys[i].ID
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, []error{err}
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	var cursor string
	if after != nil {
		cursor = string(*after)
	}

	if first == nil {
		return nil, []error{apperrors.NewInvalidDataError("missing required parameter 'first'")}
	}

	apiDefsPage, err := r.apiSvc.ListAllByPackageIDs(ctx, packageIDs, *first, cursor)
	if err != nil {
		return nil, []error{err}
	}

	err = tx.Commit()
	if err != nil {
		return nil, []error{err}
	}

	var gqlApiDefs []*graphql.APIDefinitionPage
	for _, crrApiDef := range apiDefsPage {
		apiDefs := r.apiConverter.MultipleToGraphQL(crrApiDef.Data)

		gqlApiDefs = append(gqlApiDefs, &graphql.APIDefinitionPage{Data: apiDefs, TotalCount: crrApiDef.TotalCount, PageInfo: &graphql.PageInfo{
			StartCursor: graphql.PageCursor(crrApiDef.PageInfo.StartCursor),
			EndCursor:   graphql.PageCursor(crrApiDef.PageInfo.EndCursor),
			HasNextPage: crrApiDef.PageInfo.HasNextPage,
		}})
	}

	return gqlApiDefs, nil
}

func (r *Resolver) APIDefinitionsNoPaging(ctx context.Context, obj *graphql.Package) ([]*graphql.APIDefinition, error) {
	param := dataloader.ParamApiDefNoPaging{ID: obj.ID, Ctx: ctx}
	return dataloader.ApiDefForNoPaging(ctx).ApiDefByIdNoPaging.Load(param)
}

func (r *Resolver) ApiDefinitionsDataLoaderNoPaging(keys []dataloader.ParamApiDefNoPaging) ([][]*graphql.APIDefinition, []error) {
	var ctx context.Context
	packageIDs := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		if i == 0 {
			ctx = keys[i].Ctx
			packageIDs[i] = keys[i].ID
		}
		packageIDs[i] = keys[i].ID
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, []error{err}
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	apiDefs, err := r.apiSvc.ListAllByPackageIDsNoPaging(ctx, packageIDs)
	if err != nil {
		return nil, []error{err}
	}

	err = tx.Commit()
	if err != nil {
		return nil, []error{err}
	}

	gqlApiDefs := make([][]*graphql.APIDefinition, len(packageIDs))
	for i, _ := range apiDefs {
		crrApiDefs := r.apiConverter.MultipleToGraphQL(apiDefs[i])
		gqlApiDefs[i] = crrApiDefs
	}

	return gqlApiDefs, nil
}

func (r *Resolver) EventDefinition(ctx context.Context, obj *graphql.Package, id string) (*graphql.EventDefinition, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	eventAPI, err := r.eventSvc.GetForPackage(ctx, id, obj.ID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return nil, tx.Commit()
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.eventConverter.ToGraphQL(eventAPI), nil
}

func (r *Resolver) EventDefinitions(ctx context.Context, obj *graphql.Package, group *string, first *int, after *graphql.PageCursor) (*graphql.EventDefinitionPage, error) {
	param := dataloader.ParamEventDef{ID: obj.ID, Ctx: ctx, First: first, After: after}
	return dataloader.EventDefFor(ctx).EventDefById.Load(param)
}

func (r *Resolver) EventDefinitionsDataLoader(keys []dataloader.ParamEventDef) ([]*graphql.EventDefinitionPage, []error) {
	if len(keys) == 0 {
		return nil, []error{apperrors.NewInternalError("No Packages found")}
	}

	var ctx context.Context
	var first *int
	var after *graphql.PageCursor

	packageIDs := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		if i == 0 {
			ctx = keys[i].Ctx
			first = keys[i].First
			after = keys[i].After
			packageIDs[i] = keys[i].ID
		}
		packageIDs[i] = keys[i].ID
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, []error{err}
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	var cursor string
	if after != nil {
		cursor = string(*after)
	}

	if first == nil {
		return nil, []error{apperrors.NewInvalidDataError("missing required parameter 'first'")}
	}

	eventDefsPage, err := r.eventSvc.ListAllByPackageIDs(ctx, packageIDs, *first, cursor)
	if err != nil {
		return nil, []error{err}
	}

	err = tx.Commit()
	if err != nil {
		return nil, []error{err}
	}

	var gqlEventDefs []*graphql.EventDefinitionPage
	for _, crrEventDef := range eventDefsPage {
		eventDefs := r.eventConverter.MultipleToGraphQL(crrEventDef.Data)

		gqlEventDefs = append(gqlEventDefs, &graphql.EventDefinitionPage{Data: eventDefs, TotalCount: crrEventDef.TotalCount, PageInfo: &graphql.PageInfo{
			StartCursor: graphql.PageCursor(crrEventDef.PageInfo.StartCursor),
			EndCursor:   graphql.PageCursor(crrEventDef.PageInfo.EndCursor),
			HasNextPage: crrEventDef.PageInfo.HasNextPage,
		}})
	}

	return gqlEventDefs, nil

}

func (r *Resolver) EventDefinitionsNoPaging(ctx context.Context, obj *graphql.Package) ([]*graphql.EventDefinition, error) {
	param := dataloader.ParamEventDefNoPaging{ID: obj.ID, Ctx: ctx}
	return dataloader.EventDefForNoPaging(ctx).EventDefByIdNoPaging.Load(param)
}

func (r *Resolver) EventDefinitionsDataLoaderNoPaging(keys []dataloader.ParamEventDefNoPaging) ([][]*graphql.EventDefinition, []error) {
	var ctx context.Context
	packageIDs := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		if i == 0 {
			ctx = keys[i].Ctx
			packageIDs[i] = keys[i].ID
		}
		packageIDs[i] = keys[i].ID
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, []error{err}
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	eventDefs, err := r.eventSvc.ListAllByPackageIDsNoPaging(ctx, packageIDs)
	if err != nil {
		return nil, []error{err}
	}

	err = tx.Commit()
	if err != nil {
		return nil, []error{err}
	}

	gqlEventDefs := make([][]*graphql.EventDefinition, len(packageIDs))
	for i, _ := range eventDefs {
		crrEventDefs := r.eventConverter.MultipleToGraphQL(eventDefs[i])
		gqlEventDefs[i] = crrEventDefs
	}

	return gqlEventDefs, nil
}

func (r *Resolver) Document(ctx context.Context, obj *graphql.Package, id string) (*graphql.Document, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	eventAPI, err := r.documentSvc.GetForPackage(ctx, id, obj.ID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return nil, tx.Commit()
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.documentConverter.ToGraphQL(eventAPI), nil
}

func (r *Resolver) Documents(ctx context.Context, obj *graphql.Package, first *int, after *graphql.PageCursor) (*graphql.DocumentPage, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	var cursor string
	if after != nil {
		cursor = string(*after)
	}

	if first == nil {
		return nil, apperrors.NewInvalidDataError("missing required parameter 'first'")
	}

	documentsPage, err := r.documentSvc.ListForPackage(ctx, obj.ID, *first, cursor)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlDocuments := r.documentConverter.MultipleToGraphQL(documentsPage.Data)

	return &graphql.DocumentPage{
		Data:       gqlDocuments,
		TotalCount: documentsPage.TotalCount,
		PageInfo: &graphql.PageInfo{
			StartCursor: graphql.PageCursor(documentsPage.PageInfo.StartCursor),
			EndCursor:   graphql.PageCursor(documentsPage.PageInfo.EndCursor),
			HasNextPage: documentsPage.PageInfo.HasNextPage,
		},
	}, nil
}
