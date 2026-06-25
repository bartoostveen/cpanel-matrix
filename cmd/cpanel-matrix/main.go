package main

import (
	"os"
	"time"

	env "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.bartoostveen.nl/cpanel-matrix/config"
	"go.bartoostveen.nl/cpanel-matrix/server"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	_ = env.Load()

	err, appConfig := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	server.Run(appConfig)
}
