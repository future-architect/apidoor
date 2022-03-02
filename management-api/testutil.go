package managementapi

import (
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"github.com/future-architect/apidoor/managementapi/apirouting/dynamo"
	"github.com/future-architect/apidoor/managementapi/apirouting/redis"
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

type APIDBType string

const (
	DYNAMO  APIDBType = "dynamo"
	REDIS             = "redis"
	ILLEGAL           = "illegal"
)

func GetAPIDBType(t *testing.T) APIDBType {
	switch apirouting.ApiDBDriver.(type) {
	case *dynamo.APIRouting:
		return DYNAMO
	case *redis.APIRouting:
		return REDIS
	default:
		return ILLEGAL
	}
}
