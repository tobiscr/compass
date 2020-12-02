package api

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/dataloader"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery -name=APIService -output=automock -outpkg=automock -case=underscore
type APIService interface {
	CreateInPackage(ctx context.Context, packageID string, in model.APIDefinitionInput) (string, error)
	Update(ctx context.Context, id string, in model.APIDefinitionInput) error
	Get(ctx context.Context, id string) (*model.APIDefinition, error)
	Delete(ctx context.Context, id string) error
	RefetchAPISpec(ctx context.Context, id string) (*model.APISpec, error)
	GetFetchRequest(ctx context.Context, apiDefID string) (*model.FetchRequest, error)
	ListFetchRequests(ctx context.Context, apiDefIDs []string) ([]*model.FetchRequest, error)
}

//go:generate mockery -name=RuntimeService -output=automock -outpkg=automock -case=underscore
type RuntimeService interface {
	Get(ctx context.Context, id string) (*model.Runtime, error)
}

//go:generate mockery -name=APIConverter -output=automock -outpkg=automock -case=underscore
type APIConverter interface {
	ToGraphQL(in *model.APIDefinition) *graphql.APIDefinition
	MultipleToGraphQL(in []*model.APIDefinition) []*graphql.APIDefinition
	MultipleInputFromGraphQL(in []*graphql.APIDefinitionInput) ([]*model.APIDefinitionInput, error)
	InputFromGraphQL(in *graphql.APIDefinitionInput) (*model.APIDefinitionInput, error)
	SpecToGraphQL(definitionID string, in *model.APISpec) *graphql.APISpec
}

//go:generate mockery -name=FetchRequestConverter -output=automock -outpkg=automock -case=underscore
type FetchRequestConverter interface {
	ToGraphQL(in *model.FetchRequest) (*graphql.FetchRequest, error)
	InputFromGraphQL(in *graphql.FetchRequestInput) (*model.FetchRequestInput, error)
}

//go:generate mockery -name=ApplicationService -output=automock -outpkg=automock -case=underscore
type ApplicationService interface {
	Exist(ctx context.Context, id string) (bool, error)
}

//go:generate mockery -name=PackageService -output=automock -outpkg=automock -case=underscore
type PackageService interface {
	Exist(ctx context.Context, id string) (bool, error)
}

type Resolver struct {
	transact    persistence.Transactioner
	svc         APIService
	appSvc      ApplicationService
	pkgSvc      PackageService
	rtmSvc      RuntimeService
	converter   APIConverter
	frConverter FetchRequestConverter
}

func NewResolver(transact persistence.Transactioner, svc APIService, appSvc ApplicationService, rtmSvc RuntimeService, pkgSvc PackageService, converter APIConverter, frConverter FetchRequestConverter) *Resolver {
	return &Resolver{
		transact:    transact,
		svc:         svc,
		appSvc:      appSvc,
		rtmSvc:      rtmSvc,
		pkgSvc:      pkgSvc,
		converter:   converter,
		frConverter: frConverter,
	}
}

func (r *Resolver) AddAPIDefinitionToPackage(ctx context.Context, packageID string, in graphql.APIDefinitionInput) (*graphql.APIDefinition, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	log.Infof("Adding APIDefinition to package with id %s", packageID)

	ctx = persistence.SaveToContext(ctx, tx)

	convertedIn, err := r.converter.InputFromGraphQL(&in)
	if err != nil {
		return nil, errors.Wrap(err, "while converting GraphQL input to APIDefinition")
	}

	found, err := r.pkgSvc.Exist(ctx, packageID)
	if err != nil {
		return nil, errors.Wrapf(err, "while checking existence of Package with id %s when adding APIDefinition", packageID)
	}

	if !found {
		return nil, apperrors.NewInvalidDataError("cannot add API to not existing package")
	}

	id, err := r.svc.CreateInPackage(ctx, packageID, *convertedIn)
	if err != nil {
		return nil, errors.Wrapf(err, "Error occurred while creating APIDefinition in Package with id %s", packageID)
	}

	api, err := r.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlAPI := r.converter.ToGraphQL(api)

	log.Infof("APIDefinition with id %s successfully added to Package with id %s", id, packageID)
	return gqlAPI, nil
}

func (r *Resolver) UpdateAPIDefinition(ctx context.Context, id string, in graphql.APIDefinitionInput) (*graphql.APIDefinition, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	log.Infof("Updating APIDefinition with id %s", id)

	ctx = persistence.SaveToContext(ctx, tx)

	convertedIn, err := r.converter.InputFromGraphQL(&in)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting GraphQL input to APIDefinition with id %s", id)
	}

	err = r.svc.Update(ctx, id, *convertedIn)
	if err != nil {
		return nil, err
	}

	api, err := r.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlAPI := r.converter.ToGraphQL(api)

	log.Infof("APIDefinition with id %s successfully updated.", id)
	return gqlAPI, nil
}
func (r *Resolver) DeleteAPIDefinition(ctx context.Context, id string) (*graphql.APIDefinition, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	log.Infof("Deleting APIDefinition with id %s", id)

	ctx = persistence.SaveToContext(ctx, tx)

	api, err := r.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = r.svc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	log.Infof("APIDefinition with id %s successfully deleted.", id)
	return r.converter.ToGraphQL(api), nil
}
func (r *Resolver) RefetchAPISpec(ctx context.Context, apiID string) (*graphql.APISpec, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	log.Infof("Refetching APISpec for API with id %s", apiID)

	ctx = persistence.SaveToContext(ctx, tx)

	spec, err := r.svc.RefetchAPISpec(ctx, apiID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	converted := r.converter.SpecToGraphQL(apiID, spec)
	log.Infof("Successfully refetched APISpec for APIDefinition with id %s", apiID)
	return converted, nil
}

func (r *Resolver) FetchRequest(ctx context.Context, obj *graphql.APISpec) (*graphql.FetchRequest, error) {
	params := dataloader.ParamFetchRequestApiDef{ID: obj.DefinitionID, Ctx: ctx}
	return dataloader.ForFetchRequestApiDef(ctx).FetchRequestApiDefById.Load(params)
}

func (r *Resolver) FetchRequestApiDefDataLoader(keys []dataloader.ParamFetchRequestApiDef) ([]*graphql.FetchRequest, []error) {
	if len(keys) == 0 {
		return nil, []error{apperrors.NewInternalError("No ApiDefs found")}
	}

	var ctx context.Context
	apiDefIds := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		if keys[i].ID == "" {
			return nil, []error{apperrors.NewInternalError("Cannot fetch FetchRequest. APIDefinition ID is empty")}
		}
		if i == 0 {
			ctx = keys[i].Ctx
			apiDefIds[i] = keys[i].ID
		}
		apiDefIds[i] = keys[i].ID
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, []error{err}
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	fetchRequests, err := r.svc.ListFetchRequests(ctx, apiDefIds)
	if err != nil {
		return nil, []error{err}
	}

	if fetchRequests == nil {
		return nil, nil
	}

	err = tx.Commit()
	if err != nil {
		return nil, []error{err}
	}

	var gqlFetchRequests []*graphql.FetchRequest
	for _, fr := range fetchRequests {
		fetchRequest, err := r.frConverter.ToGraphQL(fr)
		if err != nil {
			return nil, []error{err}
		}
		gqlFetchRequests = append(gqlFetchRequests, fetchRequest)
	}

	return gqlFetchRequests, nil
}
