package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type Database struct {
	Conn *sql.DB
}

func CreateDatabaseAndInitSchemaIfNotExists(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	log.Debug().
		Str("db_path", dbPath).
		Msg("Connected to database")

	err = createSchema(conn)
	if err != nil {
		return nil, err
	}

	return &Database{
		Conn: conn,
	}, nil
}

func (db *Database) Close() {
	_ = db.Conn.Close()
}
