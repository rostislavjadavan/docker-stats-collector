package main

import (
	"dsc/collector"
	"dsc/config"
	"dsc/database"
	"dsc/lib"
	"dsc/webserver"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var appConfig config.AppConfig

func main() {
	appConfig = config.LoadFromEnv()
	lib.SetupLogger(appConfig.LogLevel)

	log.Info().
		Str("log_level", appConfig.LogLevel).
		Int("interval", appConfig.Interval).
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
		webserver.CreateAndListen(appConfig, db)
	}()

	scheduler.StartBlocking()
}
