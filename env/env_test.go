package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetAppEnv(t *testing.T) {
	AppEnvFallback = "dev"
	assert.Equal(t, getAppEnv(), "dev")

	os.Setenv("RNews_ENV", "prod")
	assert.Equal(t, getAppEnv(), "prod")
}

func TestProjDir(t *testing.T) {
	assert.True(t, fileExists(ProjDir()+"/news"))
}

func fileExists(path string) bool {
	_, err := os.Open(path)
	return err == nil
}
