package main

import (
	"dsc/collector"
	"dsc/config"
	"dsc/database"
	"dsc/lib"
	"dsc/webserver"
	"embed"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"io/fs"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var appConfig config.AppConfig

//go:embed public
var embeddedFiles embed.FS

func main() {
	appConfig = config.LoadFromEnv()
	lib.SetupLogger(appConfig.LogLevel)

	log.Info().
		Str("log_level", appConfig.LogLevel).
		Msg("Starting Docker stats collection")

	db, err := database.CreateDatabaseAndInitSchemaIfNotExists(appConfig.DbPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening database")
	}
	defer db.Close()

	scheduler := gocron.NewScheduler(time.UTC)
	defer scheduler.Stop()

	err = collector.CreateAndScheduleCollector(appConfig, scheduler, db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating collector")
	}

	go func() {
		publicFS, err := fs.Sub(embeddedFiles, "public")
		if err != nil {
			log.Fatal().Err(err)
		}
		webserver.CreateAndListen(appConfig, db, publicFS)
	}()

	scheduler.StartBlocking()
}
