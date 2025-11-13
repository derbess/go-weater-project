package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Connect(config Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	//test
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

func InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS weather_data (
		id SERIAL PRIMARY KEY,
		city TEXT NOT NULL,
		temperature_c NUMERIC,
		humidity NUMERIC,
		wind_speed NUMERIC,
		feels_like NUMERIC,
		provider TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_weather_city ON weather_data(city);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating schema: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
