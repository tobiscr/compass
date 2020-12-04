package appdetails

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kyma-incubator/compass/components/connectivity-adapter/internal/appregistry/director"
	"github.com/kyma-incubator/compass/components/connectivity-adapter/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/connectivity-adapter/pkg/gqlcli"
	"github.com/kyma-incubator/compass/components/connectivity-adapter/pkg/res"
	"github.com/kyma-incubator/compass/components/connectivity-adapter/pkg/retry"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/graphqlizer"
	gcli "github.com/machinebox/graphql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const appNamePathVariable = "app-name"
const appIDPathVariable = "app-name"

//go:generate mockery -name=GraphQLRequestBuilder -output=automock -outpkg=automock -case=underscore
type GraphQLRequestBuilder interface {
	GetApplicationsByName(appName string) *gcli.Request
}

type applicationMiddleware struct {
	cliProvider gqlcli.Provider
	logger      *log.Logger
	gqlProvider graphqlizer.GqlFieldsProvider
}

func NewApplicationMiddleware(cliProvider gqlcli.Provider, logger *log.Logger) *applicationMiddleware {
	return &applicationMiddleware{cliProvider: cliProvider, logger: logger}
}

func (mw *applicationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		//appName := variables[appNamePathVariable]
		appID := variables[appIDPathVariable]

		mw.logger.Infof("resolving application with id '%s'...", appID)

		client := mw.cliProvider.GQLClient(r)
		directorCli := director.NewClient(client, &graphqlizer.Graphqlizer{}, &graphqlizer.GqlFieldsProvider{})
		//query := directorCli.GetApplicationsByNameRequest(appName)
		query := directorCli.GetApplicationByIDRequest(appID)

		var apps GqlSuccessfulAppPage
		err := retry.GQLRun(client.Run, r.Context(), query, &apps)
		if err != nil {
			wrappedErr := errors.Wrap(err, "while getting service")
			mw.logger.Error(wrappedErr)
			res.WriteError(w, wrappedErr, apperrors.CodeInternal)
			return
		}

		if len(apps.Result.Data) == 0 {
			message := fmt.Sprintf("application with id %s not found", appID)
			mw.logger.Warn(message)
			res.WriteErrorMessage(w, message, apperrors.CodeNotFound)
			return
		}

		if len(apps.Result.Data) != 1 {
			message := fmt.Sprintf("found more than 1 application with id %s", appID)
			mw.logger.Warn(message)
			res.WriteErrorMessage(w, message, apperrors.CodeInternal)
			return
		}

		app := apps.Result.Data[0]

		mw.logger.Infof("app with id '%s' details fetched successfully", appID)

		ctx := SaveToContext(r.Context(), *app)
		ctxWithCli := gqlcli.SaveToContext(ctx, client)
		requestWithCtx := r.WithContext(ctxWithCli)
		next.ServeHTTP(w, requestWithCtx)
	})
}

type GqlSuccessfulAppPage struct {
	Result graphql.ApplicationPageExt `json:"result"`
}
