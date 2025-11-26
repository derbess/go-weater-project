package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Connect(config Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test connection
	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Printf("Successfully connected to database")
	return pool, nil
}

func InitSchema(pool *pgxpool.Pool) error {
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

	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error creating schema: %w", err)
	}

	log.Printf("Database schema initialized successfully")
	return nil
}

func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
