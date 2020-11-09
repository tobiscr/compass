package integrationsystem

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
)

//go:generate mockery -name=IntegrationSystemService -output=automock -outpkg=automock -case=underscore
type IntegrationSystemService interface {
	Create(ctx context.Context, in model.IntegrationSystemInput) (string, error)
	Get(ctx context.Context, id string) (*model.IntegrationSystem, error)
	List(ctx context.Context, pageSize int, cursor string) (model.IntegrationSystemPage, error)
	Update(ctx context.Context, id string, in model.IntegrationSystemInput) error
	Delete(ctx context.Context, id string) error
}

//go:generate mockery -name=IntegrationSystemConverter -output=automock -outpkg=automock -case=underscore
type IntegrationSystemConverter interface {
	ToGraphQL(in *model.IntegrationSystem) *externalschema.IntegrationSystem
	MultipleToGraphQL(in []*model.IntegrationSystem) []*externalschema.IntegrationSystem
	InputFromGraphQL(in externalschema.IntegrationSystemInput) model.IntegrationSystemInput
}

//go:generate mockery -name=SystemAuthService -output=automock -outpkg=automock -case=underscore
type SystemAuthService interface {
	ListForObject(ctx context.Context, objectType model.SystemAuthReferenceObjectType, objectID string) ([]model.SystemAuth, error)
}

//go:generate mockery -name=SystemAuthConverter -output=automock -outpkg=automock -case=underscore
type SystemAuthConverter interface {
	ToGraphQL(in *model.SystemAuth) (*externalschema.SystemAuth, error)
}

//go:generate mockery -name=OAuth20Service -output=automock -outpkg=automock -case=underscore
type OAuth20Service interface {
	DeleteMultipleClientCredentials(ctx context.Context, auths []model.SystemAuth) error
}
type Resolver struct {
	transact persistence.Transactioner

	intSysSvc        IntegrationSystemService
	sysAuthSvc       SystemAuthService
	oAuth20Svc       OAuth20Service
	intSysConverter  IntegrationSystemConverter
	sysAuthConverter SystemAuthConverter
}

func NewResolver(transact persistence.Transactioner, intSysSvc IntegrationSystemService, sysAuthSvc SystemAuthService, oAuth20Svc OAuth20Service, intSysConverter IntegrationSystemConverter, sysAuthConverter SystemAuthConverter) *Resolver {
	return &Resolver{
		transact:         transact,
		intSysSvc:        intSysSvc,
		sysAuthSvc:       sysAuthSvc,
		oAuth20Svc:       oAuth20Svc,
		intSysConverter:  intSysConverter,
		sysAuthConverter: sysAuthConverter,
	}
}

func (r *Resolver) IntegrationSystem(ctx context.Context, id string) (*externalschema.IntegrationSystem, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	is, err := r.intSysSvc.Get(ctx, id)
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

	return r.intSysConverter.ToGraphQL(is), nil
}

func (r *Resolver) IntegrationSystems(ctx context.Context, first *int, after *externalschema.PageCursor) (*externalschema.IntegrationSystemPage, error) {
	var cursor string
	if after != nil {
		cursor = string(*after)
	}
	if first == nil {
		return nil, apperrors.NewInvalidDataError("missing required parameter 'first'")
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	intSysPage, err := r.intSysSvc.List(ctx, *first, cursor)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlIntSys := r.intSysConverter.MultipleToGraphQL(intSysPage.Data)

	return &externalschema.IntegrationSystemPage{
		Data:       gqlIntSys,
		TotalCount: intSysPage.TotalCount,
		PageInfo: &externalschema.PageInfo{
			StartCursor: externalschema.PageCursor(intSysPage.PageInfo.StartCursor),
			EndCursor:   externalschema.PageCursor(intSysPage.PageInfo.EndCursor),
			HasNextPage: intSysPage.PageInfo.HasNextPage,
		},
	}, nil
}

func (r *Resolver) RegisterIntegrationSystem(ctx context.Context, in externalschema.IntegrationSystemInput) (*externalschema.IntegrationSystem, error) {
	convertedIn := r.intSysConverter.InputFromGraphQL(in)

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	id, err := r.intSysSvc.Create(ctx, convertedIn)
	if err != nil {
		return nil, err
	}

	intSys, err := r.intSysSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlIntSys := r.intSysConverter.ToGraphQL(intSys)

	return gqlIntSys, nil
}

func (r *Resolver) UpdateIntegrationSystem(ctx context.Context, id string, in externalschema.IntegrationSystemInput) (*externalschema.IntegrationSystem, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	convertedIn := r.intSysConverter.InputFromGraphQL(in)
	err = r.intSysSvc.Update(ctx, id, convertedIn)
	if err != nil {
		return nil, err
	}

	intSys, err := r.intSysSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlIntSys := r.intSysConverter.ToGraphQL(intSys)

	return gqlIntSys, nil
}

func (r *Resolver) UnregisterIntegrationSystem(ctx context.Context, id string) (*externalschema.IntegrationSystem, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	intSys, err := r.intSysSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	auths, err := r.sysAuthSvc.ListForObject(ctx, model.IntegrationSystemReference, intSys.ID)
	if err != nil {
		return nil, err
	}

	err = r.intSysSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	err = r.oAuth20Svc.DeleteMultipleClientCredentials(ctx, auths)
	if err != nil {
		return nil, err
	}

	deletedIntSys := r.intSysConverter.ToGraphQL(intSys)

	return deletedIntSys, nil
}

func (r *Resolver) Auths(ctx context.Context, obj *externalschema.IntegrationSystem) ([]*externalschema.SystemAuth, error) {
	if obj == nil {
		return nil, apperrors.NewInternalError("Integration System cannot be empty")
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	sysAuths, err := r.sysAuthSvc.ListForObject(ctx, model.IntegrationSystemReference, obj.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	var out []*externalschema.SystemAuth
	for _, sa := range sysAuths {
		c, err := r.sysAuthConverter.ToGraphQL(&sa)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}

	return out, nil
}
