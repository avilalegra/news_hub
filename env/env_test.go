package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProjDir(t *testing.T) {
	assert.True(t, fileExists(projDir()+"/news"))
}

func fileExists(path string) bool {
	_, err := os.Open(path)
	return err == nil
}
