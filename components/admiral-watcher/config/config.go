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

package config

import (
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"reflect"
	"time"

	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
)

type Validatable interface {
	Validate() error
}

type Config struct {
	Log      *log.Config
	Database persistence.DatabaseConfig
}

func DefaultConfig() *Config {
	return &Config{
		Log: log.DefaultConfig(),
		Database: persistence.DatabaseConfig{
			User:               "director",
			Password:           "bu0ohthee0woh0equ1DahG2cu2aeceec",
			Host:               "localhost",
			Port:               "5432",
			Name:               "director",
			SSLMode:            "disable",
			MaxOpenConnections: 10,
			MaxIdleConnections: 10,
			ConnMaxLifetime:    30 * time.Minute,
		},
	}
}

func (c *Config) Validate() error {
	validatableFields := make([]Validatable, 0, 0)

	v := reflect.ValueOf(*c)
	for i := 0; i < v.NumField(); i++ {
		field, ok := v.Field(i).Interface().(Validatable)
		if ok {
			validatableFields = append(validatableFields, field)
		}
	}

	for _, item := range validatableFields {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}
