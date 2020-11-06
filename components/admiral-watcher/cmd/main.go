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

package main

import (
	"context"
	"fmt"
	"github.com/kyma-incubator/compass/components/admiral-watcher/notifications"
	"github.com/kyma-incubator/compass/components/admiral-watcher/script"
	"github.com/kyma-incubator/compass/components/admiral-watcher/templates"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/kyma-incubator/compass/components/admiral-watcher/config"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/signal"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/api"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/application"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/auth"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/document"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/eventdef"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/fetchrequest"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/integrationsystem"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/labeldef"
	mp_package "github.com/kyma-incubator/compass/components/director/internal2/domain/package"
	rt "github.com/kyma-incubator/compass/components/director/internal2/domain/runtime"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/scenarioassignment"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/version"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/webhook"
	"github.com/kyma-incubator/compass/components/director/internal2/uid"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = path.Dir(b)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	term := make(chan os.Signal)
	signal.HandleInterrupts(ctx, cancel, term)

	cfg := config.DefaultConfig()
	//err := envconfig.InitWithPrefix(&cfg, "APP")
	//fatalOnError(err)

	err := cfg.Validate()
	fatalOnError(err)

	ctx, err = log.Configure(ctx, cfg.Log)
	fatalOnError(err)

	authConverter := auth.NewConverter()
	frConverter := fetchrequest.NewConverter(authConverter)
	versionConverter := version.NewConverter()
	docConverter := document.NewConverter(frConverter)
	webhookConverter := webhook.NewConverter(authConverter)
	apiConverter := api.NewConverter(frConverter, versionConverter)
	eventAPIConverter := eventdef.NewConverter(frConverter, versionConverter)
	labelDefConverter := labeldef.NewConverter()
	labelConverter := label.NewConverter()
	intSysConverter := integrationsystem.NewConverter()
	packageConverter := mp_package.NewConverter(authConverter, apiConverter, eventAPIConverter, docConverter)
	appConverter := application.NewConverter(webhookConverter, packageConverter)
	assignmentConv := scenarioassignment.NewConverter()

	runtimeRepo := rt.NewRepository()
	applicationRepo := application.NewRepository(appConverter)
	labelRepo := label.NewRepository(labelConverter)
	labelDefRepo := labeldef.NewRepository(labelDefConverter)
	webhookRepo := webhook.NewRepository(webhookConverter)
	apiRepo := api.NewRepository(apiConverter)
	eventAPIRepo := eventdef.NewRepository(eventAPIConverter)
	docRepo := document.NewRepository(docConverter)
	fetchRequestRepo := fetchrequest.NewRepository(frConverter)
	intSysRepo := integrationsystem.NewRepository(intSysConverter)
	packageRepo := mp_package.NewRepository(packageConverter)
	scenarioAssignmentRepo := scenarioassignment.NewRepository(assignmentConv)

	uidSvc := uid.NewService()
	labelUpsertSvc := label.NewLabelUpsertService(labelRepo, labelDefRepo, uidSvc)
	scenariosSvc := labeldef.NewScenariosService(labelDefRepo, uidSvc, false)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	scenarioAssignmentEngine := scenarioassignment.NewEngine(labelUpsertSvc, labelRepo, scenarioAssignmentRepo)
	fetchRequestSvc := fetchrequest.NewService(fetchRequestRepo, httpClient, log.D().Logger)
	packageSvc := mp_package.NewService(packageRepo, apiRepo, eventAPIRepo, docRepo, fetchRequestRepo, uidSvc, fetchRequestSvc)

	runtimeSvc := rt.NewService(runtimeRepo, labelRepo, scenariosSvc, labelUpsertSvc, uidSvc, scenarioAssignmentEngine)
	appSvc := application.NewService(&DummyApplicationHideCfgProvider{}, applicationRepo, webhookRepo, runtimeRepo, labelRepo, intSysRepo, labelUpsertSvc, scenariosSvc, packageSvc, uidSvc)

	transact, closeFunc, err := persistence.Configure(log.D().Logger, cfg.Database)
	fatalOnError(err)

	defer func() {
		err := closeFunc()
		fatalOnError(err)
	}()

	runner := script.Runner{
		ScriptsLocation: fmt.Sprintf("%s/../resources", basepath),
		Resolver:        templates.Resolver{},
	}

	appLabelsHandler := &notifications.AppLabelNotificationHandler{
		RuntimeLister:      runtimeSvc,
		AppLister:          appSvc,
		AppLabelGetter:     appSvc,
		RuntimeLabelGetter: runtimeSvc,
		Transact:           transact,
		ScriptRunner:       runner,
	}

	rtLabelsHandler := &notifications.RuntimeLabelNotificationHandler{
		RuntimeGetter:  runtimeSvc,
		AppLister:      appSvc,
		AppLabelGetter: appSvc,
		Transact:       transact,
		ScriptRunner:   runner,
	}

	labelsHandler := &notifications.LabelNotificationHandler{
		Handlers: map[resource.Type]notifications.NotificationLabelHandler{
			resource.Application: appLabelsHandler,
			resource.Runtime:     rtLabelsHandler,
		},
	}

	appHandler := &notifications.AppNotificationHandler{
		ScriptRunner: runner,
	}

	rtHandler := &notifications.RtNotificationHandler{
		ScriptRunner: runner,
	}

	processor := notifications.NewNotificationProcessor(cfg.Database, map[notifications.HandlerKey]notifications.NotificationHandler{
		{
			NotificationChannel: "events",
			ResourceType:        resource.Label,
		}: labelsHandler,
		{
			NotificationChannel: "events",
			ResourceType:        resource.Application,
		}: appHandler,
		{
			NotificationChannel: "events",
			ResourceType:        resource.Runtime,
		}: rtHandler,
	})

	if err := processor.Run(ctx); err != nil {
		fatalOnError(err)
	}
}

func fatalOnError(err error) {
	if err != nil {
		log.D().Fatal(err.Error())
	}
}

type DummyApplicationHideCfgProvider struct {
}

func (d *DummyApplicationHideCfgProvider) GetApplicationHideSelectors() (map[string][]string, error) {
	return map[string][]string{}, nil
}
