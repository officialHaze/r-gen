package main

import (
	"io"
	"log"
	"os"
	"reportgenengine/api/REST/server"
	"reportgenengine/settings"
	"reportgenengine/util"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	// Generate my settings
	if err := settings.Generate(); err != nil {
		log.Fatalln(err)
	}

	// load the env file
	if err := godotenv.Load(settings.MySettings.ENV_FILE_NAME); err != nil {
		log.Fatalln(err)
	}

	// Setup lumberjack and logger
	logfile := &lumberjack.Logger{
		Filename:   "/var/log/reportgenengine.log",
		MaxSize:    100, // mb
		MaxBackups: 5,
		MaxAge:     28, // days
		Compress:   true,
	}

	// Log only in console when in dev mode
	if util.InDevMode() {
		multiwriter := io.MultiWriter(os.Stdout)
		log.SetOutput(multiwriter)
	} else {
		log.SetOutput(logfile)
	}
}

func main() {
	log.Printf(" ****** PROJECT Running In %s Environment ******", strings.ToUpper(os.Getenv("ENV")))

	util.StartTokenBucket()
	
	// Start the server
	server.Start()
}
