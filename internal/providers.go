package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"go-first-project/models"
	"net/http"
	"time"
)


func FeatchOpenWeather(ctx context.Context, client *http.Client, city string) (*models.OpenWeather, error) {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	apiKeyOpenWeather := getEnv("API_KEY_OPEN_WEATHER", "")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKeyOpenWeather)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request weather for %s: %w", city, err)
	}

	defer resp.Body.Close()

	var data models.OpenWeather

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Printf("Error decoding whether data for %s, with error %w", city, err)
		return nil, fmt.Errorf("decode weather response for %s: %w", city, err)
	}

	return &data, nil
}

func FetchWeatherApi(ctx context.Context, client *http.Client, city string) (*models.WeatherApi, error) {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	apiKeyWeatherApi := getEnv("API_KEY_WEATHER_API", "")
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&key=%s", city, apiKeyWeatherApi)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request weather for %s: %w", city, err)
	}

	defer resp.Body.Close()

	var data models.WeatherApi

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Printf("Error decoding whether data for %s, with error %w", city, err)
		return nil, fmt.Errorf("decode weather response for %s: %w", city, err)
	}

	return &data, nil

}
