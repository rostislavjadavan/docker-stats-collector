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

	_, err = scheduler.Every(appConfig.Interval).Seconds().Do(func() {
		err := statsCollector.CollectAndStoreStats()
		if err != nil {
			log.Error().Err(err).Msg("Error collecting stats")
		}
	})
	if err != nil {
		return err
	}

	log.Info().
		Int("interval", appConfig.Interval).
		Msg("Stats collector created and started")
	return nil
}
