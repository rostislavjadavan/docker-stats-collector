package database

import (
	"time"
)

type MemoryTopStat struct {
	Date          time.Time
	ContainerID   string
	ContainerName string
	CpuPercent    float64
	MemoryUsage   uint64
	MemoryLimit   uint64
}

func (d *Database) GetTopMemoryStats(image string) ([]MemoryTopStat, error) {
	query := `
	WITH daily_stats AS (
	  SELECT 
		date(timestamp) AS stat_date,
		id,
		name,
		MAX(memory_usage) AS max_memory_usage,
		memory_limit,
		cpu_percent
	  FROM container_stats
	  WHERE 
		image = ?
		AND timestamp >= date('now', '-7 days')
	  GROUP BY date(timestamp), id, name
	),
	ranked_stats AS (
	  SELECT 
		*,
		ROW_NUMBER() OVER (PARTITION BY stat_date ORDER BY max_memory_usage DESC) AS rank
	  FROM daily_stats
	)
	SELECT 
	  stat_date,
	  id,
	  name,
	  max_memory_usage,
	  memory_limit,
	  cpu_percent
	FROM ranked_stats
	WHERE rank <= 5
	ORDER BY stat_date DESC, max_memory_usage DESC;
	`

	rows, err := d.Conn.Query(query, image)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []MemoryTopStat
	for rows.Next() {
		var stat MemoryTopStat
		var dateStr string
		err := rows.Scan(&dateStr, &stat.ContainerID, &stat.ContainerName, &stat.MemoryUsage, &stat.MemoryLimit, &stat.CpuPercent)
		if err != nil {
			return nil, err
		}
		stat.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

type CpuTopStat struct {
	Date          time.Time
	ContainerID   string
	ContainerName string
	CpuPercent    float64
	MemoryUsage   uint64
	MemoryLimit   uint64
}

func (d *Database) GetTopCpuStats(image string) ([]CpuTopStat, error) {
	query := `
	WITH daily_stats AS (
	  SELECT 
		date(timestamp) AS stat_date,
		id,
		name,
		MAX(cpu_percent) AS max_cpu_percent,
		memory_usage,
		memory_limit
	  FROM container_stats
	  WHERE 
		image = ?
		AND timestamp >= date('now', '-7 days')
	  GROUP BY date(timestamp), id, name
	),
	ranked_stats AS (
	  SELECT 
		*,
		ROW_NUMBER() OVER (PARTITION BY stat_date ORDER BY max_cpu_percent DESC) AS rank
	  FROM daily_stats
	)
	SELECT 
	  stat_date,
	  id,
	  name,
	  max_cpu_percent,
	  memory_usage,
	  memory_limit
	FROM ranked_stats
	WHERE rank <= 5
	ORDER BY stat_date DESC, max_cpu_percent DESC;
	`

	rows, err := d.Conn.Query(query, image)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []CpuTopStat
	for rows.Next() {
		var stat CpuTopStat
		var dateStr string
		err := rows.Scan(&dateStr, &stat.ContainerID, &stat.ContainerName, &stat.CpuPercent, &stat.MemoryUsage, &stat.MemoryLimit)
		if err != nil {
			return nil, err
		}
		stat.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
