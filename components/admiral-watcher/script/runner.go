package script

import (
	"context"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/admiral-watcher/templates"
	"github.com/kyma-incubator/compass/components/admiral-watcher/types"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

type Runner struct {
	ScriptsLocation string
	Resolver        templates.Resolver
}

func (r *Runner) Run(ctx context.Context, scriptName string, args ...string) error {
	cmdArgs := append([]string{scriptName}, args...)
	log.C(ctx).Infof("Executing %s %s", scriptName, strings.Join(args, " "))

	cmd := exec.Command("/bin/sh", cmdArgs...)
	outputBytes, err := cmd.Output()
	if err != nil {
		return errors.Errorf("failed executing script %s: %s", scriptName, err)
	}

	log.C(ctx).Infof("\n=====================\n%s=====================\n", string(outputBytes))

	return nil
}

func (r *Runner) ApplyDependency(ctx context.Context, dep types.Dependency, remote string) error {
	dependencyString := r.Resolver.ResolveDependency(dep)
	return r.Run(ctx, "dependency_applier.sh", dependencyString, remote)
}

func (r *Runner) DeleteDependency(ctx context.Context, dep string, remote string) error {
	return r.Run(ctx, "dependency_cleaner.sh", dep, remote)
}
