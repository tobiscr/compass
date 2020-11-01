package templates

import (
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/utils"
	"github.com/kyma-incubator/compass/components/admiral-watcher/types"
)

type Resolver struct {
}

func (r *Resolver) ResolveDependency(dependency types.Dependency) string {
	//TODO templates as optimization?
	return utils.ToYAML(dependency)
}
