package database

import "database/sql"

func createSchema(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS container_stats (
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
		return err
	}

	return nil
}
