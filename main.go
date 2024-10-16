package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	_ "github.com/mattn/go-sqlite3"
)

var (
	interval int
	dbPath   string
)

func init() {
	flag.IntVar(&interval, "interval", 5, "Data collection interval in seconds")
	flag.StringVar(&dbPath, "db", "container_stats.db", "Path to SQLite database file")
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	log.Printf("Starting Docker stats collection. Interval: %d seconds, DB: %s", interval, dbPath)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.46"))
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS container_stats (
		id TEXT,
		name TEXT,
		image TEXT,
		timestamp DATETIME,
		cpu_percent REAL,
		memory_usage INTEGER,
		memory_limit INTEGER,
		network_rx INTEGER,
		network_tx INTEGER
	)`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	log.Println("Data collection started")
	for range ticker.C {
		collectAndStoreStats(cli, db)
	}
}

func collectAndStoreStats(cli *client.Client, db *sql.DB) {
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		log.Printf("Error listing containers: %v", err)
		return
	}

	stmt, err := db.Prepare(`INSERT INTO container_stats
		(id, name, image, timestamp, cpu_percent, memory_usage, memory_limit, network_rx, network_tx)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return
	}
	defer stmt.Close()

	for _, containerInfo := range containers {
		stats, err := cli.ContainerStats(context.Background(), containerInfo.ID, false)
		if err != nil {
			log.Printf("Error fetching stats for container %s: %v", containerInfo.ID[:12], err)
			continue
		}
		defer stats.Body.Close()

		var statsJSON container.StatsResponse
		err = json.NewDecoder(stats.Body).Decode(&statsJSON)
		if err != nil {
			log.Printf("Error decoding stats for container %s: %v", containerInfo.ID[:12], err)
			continue
		}

		cpuDelta := float64(statsJSON.CPUStats.CPUUsage.TotalUsage - statsJSON.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(statsJSON.CPUStats.SystemUsage - statsJSON.PreCPUStats.SystemUsage)
		cpuPercent := 0.0
		if systemDelta > 0 && cpuDelta > 0 {
			cpuPercent = (cpuDelta / systemDelta) * float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage)) * 100.0
		}

		_, err = stmt.Exec(
			containerInfo.ID[:12],
			strings.Join(containerInfo.Names, "/"),
			containerInfo.Image,
			time.Now().Format(time.RFC3339),
			cpuPercent,
			statsJSON.MemoryStats.Usage,
			statsJSON.MemoryStats.Limit,
			statsJSON.Networks["eth0"].RxBytes,
			statsJSON.Networks["eth0"].TxBytes,
		)
		if err != nil {
			log.Printf("Error inserting data for container %s: %v", containerInfo.ID[:12], err)
			continue
		}

		// log.Printf("Collected stats for container %s (Image: %s)", strings.Join(containerInfo.Names, "/"), containerInfo.Image)
	}
}
