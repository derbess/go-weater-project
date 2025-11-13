package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-first-project/models"
)

type WeatherRepository struct {
	db *sql.DB
}


func NewWeatherRepository(db *sql.DB) *WeatherRepository {
	return &WeatherRepository{db: db}
}

func (r *WeatherRepository) CreateWeather(ctx context.Context, weather *models.WeatherData) error {

	query := `
		INSERT INTO weather_data (city, temperature_c, humidity, wind_speed, feels_like, provider)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		weather.City,
		weather.Temp, 
		weather.Humidity, 
		weather.Wind,
		weather.FeelsLike, 
		weather.Provider,
	)

	if err != nil {
		return fmt.Errorf("failed to insert weather: %w", err)
	}
	
	return nil
}