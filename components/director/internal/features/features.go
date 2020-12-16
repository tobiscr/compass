package features

import "github.com/kyma-incubator/compass/components/director/pkg/scenario"

type Config struct {
	DefaultScenarioEnabled bool `envconfig:"default=true,APP_DEFAULT_SCENARIO_ENABLED"`
	scenario.CallbackConfig
}
