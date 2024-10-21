package collector

import (
	"context"
	"database/sql"
	"dsc/config"
	"dsc/database"
	"encoding/json"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

var Resolution = time.Minute

type StatsCollector struct {
	cli             *client.Client
	insertStatement *sql.Stmt
}

func NewStatsCollector(config config.AppConfig, db *database.Database) (*StatsCollector, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(config.DockerClientVersion))
	if err != nil {
		return nil, err
	}
	log.Debug().
		Str("version", config.DockerClientVersion).
		Msg("Docker client created")

	stmt, err := db.Conn.Prepare(`INSERT INTO container_stats
		(id, name, image, tag, timestamp, cpu_percent, memory_usage, memory_limit, network_rx, network_tx)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	return &StatsCollector{
		cli:             dockerClient,
		insertStatement: stmt,
	}, nil
}

func (s *StatsCollector) CollectAndStoreStats() error {
	containers, err := s.cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return err
	}

	for _, containerInfo := range containers {
		stats, err := s.cli.ContainerStats(context.Background(), containerInfo.ID, false)
		if err != nil {
			log.
				Error().
				Str("containerId", containerInfo.ID[:12]).
				Str("containerName", namesToString(containerInfo.Names)).
				Err(err).
				Msg("Error fetching stats for container")
			continue
		}
		defer stats.Body.Close()

		var statsJSON container.StatsResponse
		err = json.NewDecoder(stats.Body).Decode(&statsJSON)
		if err != nil {
			log.
				Error().
				Str("containerId", containerInfo.ID[:12]).
				Str("containerName", namesToString(containerInfo.Names)).
				Err(err).
				Msg("Error decoding stats for container")
			continue
		}

		cpuDelta := float64(statsJSON.CPUStats.CPUUsage.TotalUsage - statsJSON.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(statsJSON.CPUStats.SystemUsage - statsJSON.PreCPUStats.SystemUsage)
		numCPUs := float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage))
		cpuPercent := 0.0
		if systemDelta > 0 && cpuDelta > 0 {
			cpuPercent = (cpuDelta / systemDelta) * numCPUs * 100.0
		}

		imageInfo := parseDockerImage(containerInfo.Image)

		_, err = s.insertStatement.Exec(
			containerInfo.ID[:12],
			namesToString(containerInfo.Names),
			imageInfo.BaseImage,
			imageInfo.Tag,
			time.Now().Format(time.RFC3339),
			cpuPercent,
			statsJSON.MemoryStats.Usage,
			statsJSON.MemoryStats.Limit,
			statsJSON.Networks["eth0"].RxBytes,
			statsJSON.Networks["eth0"].TxBytes,
		)
		if err != nil {
			log.
				Error().
				Str("containerId", containerInfo.ID[:12]).
				Str("containerName", namesToString(containerInfo.Names)).
				Err(err).
				Msg("Error inserting data for container")
			continue
		}
		log.Debug().
			Str("containerId", containerInfo.ID[:12]).
			Str("containerName", namesToString(containerInfo.Names)).
			Msg("Data collected for container")
	}
	return nil
}

func (s *StatsCollector) Close() {
	_ = s.cli.Close()
	_ = s.insertStatement.Close()
}

func namesToString(names []string) string {
	return strings.TrimLeft(strings.Join(names, "|"), "/")
}
