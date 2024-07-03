package storedb

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"todo/config"

	_ "github.com/mattn/go-sqlite3"
)

type SQLStorage struct {
	db *sql.DB
}

func NewSQLStorage(dbPath string) *SQLStorage {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Подключено к SQLite...")

	return &SQLStorage{db: db}
}

func (s *SQLStorage) Init() (*sql.DB, error) {
	if err := s.createTasksTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *SQLStorage) createTasksTable() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date VARCHAR(8) NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat VARCHAR(128) CHECK(length(repeat) <= 128)
	);`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func GetDBFilePath() string {
	dbFilePath := config.Envs.DBfile
	if dbFilePath == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dbFilePath = filepath.Join(currentDir, config.Envs.DBfile)
	}
	return dbFilePath
}
