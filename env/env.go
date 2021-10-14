package env

import (
	"github.com/joho/godotenv"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

const APP_ENV = "dev"

func init() {
	envDir := projDir() + "/env"
	env := os.Getenv("APP_ENV")
	if "" == env {
		env = APP_ENV
	}

	if err := godotenv.Load(envDir + "/.env." + env + ".local"); err != nil {
		panic(err)
	}
	if 	err := godotenv.Load(envDir + "/.env." + env); err != nil {
		panic(err)
	}
	if 	err := godotenv.Load(envDir + "/.env"); err != nil {
		panic(err)
	}
}

func projDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
