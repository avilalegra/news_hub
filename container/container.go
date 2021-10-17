package container

import (
	"avilego.me/news_hub/env"

	"log"
	"os"
)

func GetLogger() *log.Logger {
	file, _ := os.OpenFile(env.LogFile(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	return log.New(file, "", log.LstdFlags)
}
