package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var AppEnvFallback = "test"

func init() {
	envDir := ProjDir() + "/env"
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

func ProjDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func LogFile() string {
	return fmt.Sprintf("%s/log/%s.log", ProjDir(), getAppEnv())
}
