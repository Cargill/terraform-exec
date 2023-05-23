package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateRepProviderCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateRepProviderCmd, err := tf.stateRepProviderCmd(context.Background(), "testsource", "testdestination")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"replace-provider",
			"-no-color",
			"-auto-approve",
			"-lock-timeout=0s",
			"-lock=true",
			"--",
			"testsource",
			"testdestination",
		}, nil, stateRepProviderCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateRepProviderCmd, err := tf.stateRepProviderCmd(context.Background(), "testsrc", "testdest", Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"replace-provide",
			"-no-color",
			"-auto-approve",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-lock=false",
			"--",
			"testsrc",
			"testdest",
		}, nil, stateRepProviderCmd)
	})
}
