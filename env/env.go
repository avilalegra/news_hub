package env

import (
	"github.com/joho/godotenv"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var AppEnvFallback = "test"

func init() {
	envDir := projDir() + "/env"
	env := getAppEnv()

	if err := godotenv.Load(envDir + "/.env." + env + ".local"); err != nil {
		panic(err)
	}
	if err := godotenv.Load(envDir + "/.env." + env); err != nil {
		panic(err)
	}
	if err := godotenv.Load(envDir + "/.env"); err != nil {
		panic(err)
	}
}

func getAppEnv() string {
	env := os.Getenv("APP_ENV_FALLBACK")
	if "" == env {
		env = AppEnvFallback
	}
	return env
}

func projDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
