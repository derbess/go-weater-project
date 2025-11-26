package repository

import (
	"context"
	"fmt"
	"go-first-project/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WeatherRepository struct {
	pool *pgxpool.Pool
}

func NewWeatherRepository(pool *pgxpool.Pool) *WeatherRepository {
	return &WeatherRepository{pool: pool}
}

func (r *WeatherRepository) CreateWeather(ctx context.Context, weather *models.WeatherData) error {

	query := `
		INSERT INTO weather_data (city, temperature_c, humidity, wind_speed, feels_like, provider)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(
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
