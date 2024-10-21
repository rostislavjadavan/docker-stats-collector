package collector

import (
	"dsc/config"
	"dsc/database"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

func CreateAndScheduleCollector(appConfig config.AppConfig, scheduler *gocron.Scheduler, db *database.Database) error {
	statsCollector, err := NewStatsCollector(appConfig, db)
	if err != nil {
		return err
	}

	_, err = scheduler.Every(Resolution).Do(func() {
		err := statsCollector.CollectAndStoreStats()
		if err != nil {
			log.Error().Err(err).Msg("Error collecting stats")
		}
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("interval", Resolution.String()).
		Msg("Stats collector created and started")

	_, err = scheduler.Every(CleanUpResolution).Do(func() {
		rowsAffected, err := db.CleanupOldData()
		if err != nil {
			log.Error().Err(err).Msg("Error cleanup old stats")
		}
		log.Info().
			Int64("rows_affected", rowsAffected).
			Msg("Cleaning up old data finished")
	})
	if err != nil {
		return err
	}

	return nil
}
