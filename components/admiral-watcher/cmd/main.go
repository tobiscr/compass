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
	"github.com/kyma-incubator/compass/components/admiral-watcher/config"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/env"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/signal"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/uuid"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	term := make(chan os.Signal)
	signal.HandleInterrupts(ctx, cancel, term)

	env, err := env.Default(ctx, config.AddPFlags)
	fatalOnError(err)

	cfg, err := config.New(env)
	fatalOnError(err)

	err = cfg.Validate()
	fatalOnError(err)

	ctx, err = log.Configure(ctx, cfg.Log)
	fatalOnError(err)

	uuidSrv := uuid.NewService()

}

func fatalOnError(err error) {
	if err != nil {
		log.D().Fatal(err.Error())
	}
}
