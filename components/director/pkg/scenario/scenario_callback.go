/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package scenario

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/internal/domain/api"
	"github.com/kyma-incubator/compass/components/director/internal/domain/application"
	"github.com/kyma-incubator/compass/components/director/internal/domain/auth"
	"github.com/kyma-incubator/compass/components/director/internal/domain/document"
	"github.com/kyma-incubator/compass/components/director/internal/domain/eventdef"
	"github.com/kyma-incubator/compass/components/director/internal/domain/fetchrequest"
	packageutil "github.com/kyma-incubator/compass/components/director/internal/domain/package"
	"github.com/kyma-incubator/compass/components/director/internal/domain/runtime"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/domain/version"
	"github.com/kyma-incubator/compass/components/director/internal/domain/webhook"
	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/tidwall/sjson"
	"net/http"
)

const defaultName = "DEFAULT"

var (
	httpDestinationType          = "HTTP"
	noAuthnDestination           = "NoAuthentication"
	basicAuthnDestination        = "NoAuthentication"
	internetDestinationProxyType = "Internet"
)

type CallbackConfig struct {
	Client       Client
	TenantName   string `envconfig:"APP_CALLBACK_CONFIG_TENANT_NAME"`
	RuntimeName  string `envconfig:"APP_CALLBACK_CONFIG_RUNTIME_NAME"`
	ScenarioName string `envconfig:"APP_CALLBACK_CONFIG_SCENARIO_NAME"`
}

type Client struct {
	ClientID     string `envconfig:"APP_CALLBACK_CONFIG_CLIENT_ID"`
	ClientSecret string `envconfig:"APP_CALLBACK_CONFIG_CLIENT_SECRET"`
	TokenURL     string `envconfig:"APP_CALLBACK_CONFIG_TOKEN_URL"`
	CallbackURL  string `envconfig:"APP_CALLBACK_CONFIG_CALLBACK_URL"`
}

type Destination struct {
	Name                     *string `json:"Name,omitempty"`
	Type                     *string `json:"Type,omitempty"`
	URL                      *string `json:"URL,omitempty"`
	Authentication           *string `json:"Authentication,omitempty"`
	ProxyType                *string `json:"ProxyType,omitempty"`
	User                     *string `json:"User,omitempty"`
	Password                 *string `json:"Password,omitempty"`
	VerificationKeys         *string `json:"VerificationKeys,omitempty"`
	CloudConnectorLocationId *string `json:"CloudConnectorLocationId,omitempty"`
}

type Callback struct {
	client       *http.Client
	transact     persistence.Transactioner
	appRepo      application.ApplicationRepository
	runtimeRepo  runtime.RuntimeRepository
	packageRepo  packageutil.PackageRepository
	apiRepo      api.APIRepository
	tenantRepo   tenant.TenantMappingRepository
	callbackURL  string
	tenantName   string
	runtimeName  string
	scenarioName string
}

func NewCallbackDirective(transact persistence.Transactioner, client *http.Client, config *CallbackConfig) Callback {
	authConverter := auth.NewConverter()

	frConverter := fetchrequest.NewConverter(authConverter)
	versionConverter := version.NewConverter()
	docConverter := document.NewConverter(frConverter)

	apiConverter := api.NewConverter(frConverter, versionConverter)
	eventAPIConverter := eventdef.NewConverter(frConverter, versionConverter)

	webhookConverter := webhook.NewConverter(authConverter)
	packageConverter := packageutil.NewConverter(authConverter, apiConverter, eventAPIConverter, docConverter)
	appConverter := application.NewConverter(webhookConverter, packageConverter)

	tenantConverter := tenant.NewConverter()

	return Callback{
		client:       client,
		transact:     transact,
		appRepo:      application.NewRepository(appConverter),
		runtimeRepo:  runtime.NewRepository(),
		packageRepo:  packageutil.NewRepository(packageConverter),
		apiRepo:      api.NewRepository(apiConverter),
		tenantRepo:   tenant.NewRepository(tenantConverter),
		callbackURL:  config.Client.CallbackURL,
		tenantName:   config.TenantName,
		runtimeName:  config.RuntimeName,
		scenarioName: config.ScenarioName,
	}
}

func (c Callback) ScenarioCallback(ctx context.Context, _ interface{}, next graphql.Resolver) (res interface{}, err error) {
	tx, err := c.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while opening the db transaction.")
		return nil, err
	}
	defer c.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	tenant, err := c.tenantRepo.GetByExternalTenantName(ctx, c.tenantName)
	if err != nil {
		return nil, err
	}

	tenantID := tenant.ID

	// get apps in scenario before operation
	log.C(ctx).Infof("Fetching all applications for %q scenario BEFORE operation", c.scenarioName)
	appsBefore, err := c.fetchAppsInScenario(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	log.C(ctx).Infof("Found %s applications for %q scenario BEFORE operation", len(appsBefore), c.scenarioName)

	log.C(ctx).Info("Proceeding with operation...")
	res, err = next(ctx)
	if err != nil {
		return res, err
	}

	// get apps in scenario after operation
	log.C(ctx).Infof("Fetching all applications for %q scenario AFTER operation", c.scenarioName)
	appsAfter, err := c.fetchAppsInScenario(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	log.C(ctx).Infof("Found %s applications for %q scenario AFTER operation", len(appsAfter), c.scenarioName)

	// diff
	removedApps := make([]string, 0)
	for appName := range appsBefore {
		_, exists := appsAfter[appName]
		if !exists {
			removedApps = append(removedApps, appName)
		}
	}
	log.C(ctx).Info("Found %s removed applications from scenario %q: %+v", len(removedApps), c.scenarioName, removedApps)

	newApps := make([]model.Application, 0)
	for appName, app := range appsAfter {
		_, exists := appsBefore[appName]
		if !exists {
			newApps = append(newApps, app)
		}
	}
	log.C(ctx).Info("Found %s new applications from scenario %q: %+v", len(newApps), c.scenarioName, newApps)

	// deleted apps -> DELETE destination
	for _, appName := range removedApps {
		url := fmt.Sprintf("%s/%s", c.callbackURL, appName)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			return nil, apperrors.NewInternalError("unable to construct deletion request:", err.Error())
		}

		log.C(ctx).Info("Proceeding with DELETE %s request...", url)
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, apperrors.NewInternalError("unable to fetch destinations:", err.Error())
		}

		if resp.StatusCode != http.StatusOK {
			log.C(ctx).Infof("Retrieved response status %s when trying to delete destination %q", resp.StatusCode, appName)
		} else {
			log.C(ctx).Infof("Successfully deleted destination %q", appName)
		}
	}

	// created apps -> CREATE destination
	for _, app := range newApps {
		destination := Destination{
			Name:           &app.Name,
			Type:           &httpDestinationType,
			ProxyType:      &internetDestinationProxyType,
			Authentication: &noAuthnDestination,
		}

		packages, err := c.packageRepo.ListByApplicationID(ctx, tenantID, app.ID, 100, "")
		if err != nil {
			return nil, err
		}

		var defaultPackage *model.Package
		for _, pkg := range packages.Data {
			if pkg.Name == defaultName {
				defaultPackage = pkg
				break
			}
		}

		headers := make(map[string][]string, 0)
		if defaultPackage != nil {
			defaultCredentials := defaultPackage.DefaultInstanceAuth
			if defaultCredentials != nil {
				if defaultCredentials.Credential.Basic != nil {
					destination.Authentication = &basicAuthnDestination
					destination.User = &defaultCredentials.Credential.Basic.Username
					destination.Password = &defaultCredentials.Credential.Basic.Password
				}

				headers = defaultCredentials.AdditionalHeaders
			}

			apis, err := c.apiRepo.ListForPackage(ctx, tenantID, defaultPackage.ID, 100, "")
			if err != nil {
				return nil, err
			}

			var defaultAPI *model.APIDefinition
			for _, api := range apis.Data {
				if api.Name == defaultName {
					defaultAPI = api
					break
				}
			}

			if defaultAPI != nil {
				destination.URL = &defaultAPI.TargetURL

				if defaultAPI.Group != nil {
					destination.ProxyType = defaultAPI.Group
				}
			}
		}

		body, err := json.Marshal(&destination)
		if err != nil {
			return nil, err
		}

		for key, values := range headers {
			body, err = sjson.SetBytes(body, key, values[0])
			if err != nil {
				return nil, err
			}
		}

		log.C(ctx).Info("Proceeding with POST %s request, body = %s", c.callbackURL, string(body))
		req, err := http.NewRequest(http.MethodPost, c.callbackURL, bytes.NewBuffer(body))
		if err != nil {
			return nil, apperrors.NewInternalError("unable to construct deletion request:", err.Error())
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, apperrors.NewInternalError("unable to fetch destinations:", err.Error())
		}

		if resp.StatusCode != http.StatusCreated {
			log.C(ctx).Infof("retrieved response status %s when trying to create destination %q", resp.StatusCode, app.Name)
		} else {
			log.C(ctx).Infof("Successfully created destination %q", app.Name)
		}
	}

	return res, err
}

func (c Callback) fetchAppsInScenario(ctx context.Context, tenantID string) (map[string]model.Application, error) {
	query := fmt.Sprintf(`$[*] ? (@ == "%s")`, c.scenarioName)
	runtimes, err := c.runtimeRepo.List(ctx, tenantID, []*labelfilter.LabelFilter{{Key: "scenarios", Query: &query}}, 100, "")
	if err != nil {
		return nil, apperrors.NewInternalError("unable to fetch runtimes for tenant with ID = ", tenantID)
	}

	runtimeInScenario := false
	for _, runtime := range runtimes.Data {
		if runtime.Name == c.runtimeName {
			runtimeInScenario = true
			break
		}
	}

	if !runtimeInScenario {
		return map[string]model.Application{}, nil
	}

	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, apperrors.NewInvalidDataError("tenantID is not UUID")
	}

	apps, err := c.appRepo.ListByScenarios(ctx, tenantUUID, []string{c.scenarioName}, 100, "", nil)
	if err != nil {
		return nil, apperrors.NewInternalError("unable to fetch applications for tenant with ID = ", tenantID)
	}

	appsMap := make(map[string]model.Application, 0)
	for i, app := range apps.Data {
		appsMap[app.Name] = *apps.Data[i]
	}

	return appsMap, nil
}
