package env

import (
	"os"
	"strings"
)

func RunsInK8S() bool {
	return strings.EqualFold(os.Getenv("LogEnv"), "k8s")
}
