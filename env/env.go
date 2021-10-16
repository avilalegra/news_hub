package env

import (
	"github.com/joho/godotenv"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var AppEnvFallback = "dev"

func init() {
	envDir := projDir() + "/env"
	env := getAppEnv()

	godotenv.Load(envDir + "/.env." + env + ".local")
	godotenv.Load(envDir + "/.env." + env)
	godotenv.Load(envDir + "/.env.local")
	godotenv.Load(envDir + "/.env")
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
