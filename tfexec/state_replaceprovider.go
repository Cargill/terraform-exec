package tfexec

import (
	"context"
	"os/exec"
	"strconv"
)

type stateRepProviderConfig struct {
	backup      string
	lock        bool
	lockTimeout string
	state       string
	stateOut    string
}

var defaultStateRepProviderOptions = stateRepProviderConfig{
	lock:        true,
	lockTimeout: "0s",
}

// stateRepProviderOption represents options used in the Refresh method.
type stateRepProviderOption interface {
	configureStateRepProvider(*stateRepProviderConfig)
}

func (opt *BackupOption) configureStateRepProvider(conf *stateRepProviderConfig) {
	conf.backup = opt.path
}

func (opt *LockOption) configureStateRepProvider(conf *stateRepProviderConfig) {
	conf.lock = opt.lock
}

func (opt *LockTimeoutOption) configureStateRepProvider(conf *stateRepProviderConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *StateOption) configureStateRepProvider(conf *stateRepProviderConfig) {
	conf.state = opt.path
}

func (opt *StateOutOption) configureStateRepProvider(conf *stateRepProviderConfig) {
	conf.stateOut = opt.path
}

// StateRepProvider represents the terraform state replace-provider subcommand.
func (tf *Terraform) StateRepProvider(ctx context.Context, source string, destination string, opts ...stateRepProviderOption) error {
	cmd, err := tf.stateRepProviderCmd(ctx, source, destination, opts...)
	if err != nil {
		return err
	}
	return tf.runTerraformCmd(ctx, cmd)
}

func (tf *Terraform) stateRepProviderCmd(ctx context.Context, source string, destination string, opts ...stateRepProviderOption) (*exec.Cmd, error) {
	c := defaultStateRepProviderOptions

	for _, o := range opts {
		o.configureStateRepProvider(&c)
	}

	args := []string{"state", "replace-provider", "-no-color", "-auto-approve", "--"}

	// string opts: only pass if set
	if c.backup != "" {
		args = append(args, "-backup="+c.backup)
	}
	if c.lockTimeout != "" {
		args = append(args, "-lock-timeout="+c.lockTimeout)
	}
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}
	if c.stateOut != "" {
		args = append(args, "-state-out="+c.stateOut)
	}

	// boolean and numerical opts: always pass
	args = append(args, "-lock="+strconv.FormatBool(c.lock))

	// positional arguments
	args = append(args, source)
	args = append(args, destination)

	return tf.buildTerraformCmd(ctx, nil, args...), nil
}
