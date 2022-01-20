package gateway

import (
	"os/exec"
	"strings"
	"testing"
)

func Setup(t *testing.T, commands ...string) {
	t.Helper()
	for _, cmd := range commands {
		args := strings.Split(cmd, " ")
		if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
			t.Errorf("setup error: %v %s %s", err, out, cmd)
		}
	}
}

func Teardown(t *testing.T, commands ...string) {
	t.Helper()
	for _, cmd := range commands {
		args := strings.Split(cmd, " ")
		if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
			t.Errorf("teardown error: %v %s %s", err, out, cmd)
		}
	}
}
