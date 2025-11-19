package services

import (
	"context"
	"fmt"
	"go-first-project/client"
	"go-first-project/models"
	"go-first-project/repository"
	"log"
	"time"
)

type WeatcherService struct {
	repo   *repository.WeatherRepository
	client *client.Client
}

func NewWeatcherService(repo *repository.WeatherRepository, client *client.Client) *WeatcherService {
	return &WeatcherService{
		repo:   repo,
		client: client,
	}
}

func (s *WeatcherService) GetWeather(ctx context.Context, city string) (*models.WeatherApiResponse, error) {
	select {
	case <-ctx.Done():
		log.Printf("GetWeather: context cancelled before starting")
		return nil, ctx.Err()
	default:
	}

	start := time.Now()

	apiCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	openWeatherChan := make(chan *models.OpenWeather, 1)
	weatherApiChan := make(chan *models.WeatherApi, 1)
	errChan := make(chan error, 2)

	go func() {
		select {
		case <-apiCtx.Done():
			log.Printf("OpenWeather goroutine: context cancelled")
			return
		default:
		}

		w, err := s.client.FetchOpenWeather(apiCtx, city)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("OpenWeather cancelled: %v\n", ctx.Err())
			} else {
				log.Printf("Error fetching OpenWeather: %v\n", err)
			}
			errChan <- err
			return
		}
		log.Printf("OpenWeather response time: %v\n", time.Since(start))
		openWeatherChan <- w
	}()

	go func() {
		select {
		case <-apiCtx.Done():
			log.Printf("WeatherAPI goroutine: context cancelled")
			return
		default:
		}

		w2, err := s.client.FetchWeatherApi(apiCtx, city)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("WeatherAPI cancelled: %v\n", ctx.Err())
			} else {
				log.Printf("Error fetching WeatherApi: %v\n", err)
			}
			errChan <- err
			return
		}
		log.Printf("WeatherAPI response time: %v\n", time.Since(start))
		weatherApiChan <- w2
	}()

	var response models.WeatherApiResponse

	select {
	case w := <-openWeatherChan:
		cancel()
		response.Temp = w.Main.Temp
		response.FeelsLike = w.Main.FeelsLike
		response.Wind = w.Wind.Speed
		response.City = w.Name
		response.Humidity = w.Main.Humidity
		response.Provider = "OpenWeather"
		log.Printf("Using OpenWeather response for %s", city)

	case w2 := <-weatherApiChan:
		cancel()
		response.Temp = w2.Current.Temp
		response.FeelsLike = w2.Current.FeelsLike
		response.Wind = w2.Current.Wind
		response.City = w2.Location.Name
		response.Humidity = w2.Current.Humidity
		response.Provider = "WeatherApi"
		log.Printf("Using WeatherAPI response for %s", city)

	case <-apiCtx.Done():
		if ctx.Err() == context.Canceled {
			log.Printf("Weather fetch cancelled (shutdown) for city: %s", city)
			return nil, ctx.Err()
		}
		log.Printf("Timeout - no response received for city: %s", city)
		return nil, fmt.Errorf("timeout - no response received for city: %s", city)
	}

	select {
	case <-ctx.Done():
		log.Printf("Skipping database save: context cancelled")
		return &response, nil
	default:
	}

	weatherData := &models.WeatherData{
		City:      response.City,
		Temp:      response.Temp,
		Humidity:  response.Humidity,
		Wind:      response.Wind,
		FeelsLike: response.FeelsLike,
		Provider:  response.Provider,
	}

	err := s.repo.CreateWeather(ctx, weatherData)
	if err != nil {
		if ctx.Err() != nil {
			log.Printf("Database save cancelled: %v\n", ctx.Err())
		} else {
			log.Printf("Warning: failed to save to DB: %v\n", err)
		}
	} else {
		log.Printf("Successfully saved weather data for %s to database", city)
	}

	return &response, nil
}
