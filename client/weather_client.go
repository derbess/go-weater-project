package client

import (
	"context"
	"encoding/json"
	"fmt"
	"go-first-project/models"
	"net/http"
	"time"
)

type Client struct {
	cli               *http.Client
	openWeatherApiKey string
	weatherApiKey     string
}

func NewClient(openWeatherApiKey, weatherApiKey string) *Client {
	return &Client{
		cli: &http.Client{
			Timeout: 5 * time.Second,
		},
		openWeatherApiKey: openWeatherApiKey,
		weatherApiKey:     weatherApiKey,
	}
}

func (c *Client) FetchOpenWeather(ctx context.Context, city string) (*models.OpenWeather, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
		city, c.openWeatherApiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request weather for %s: %w", city, err)
	}
	defer resp.Body.Close()

	var data models.OpenWeather
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode weather response for %s: %w", city, err)
	}

	return &data, nil
}

func (c *Client) FetchWeatherApi(ctx context.Context, city string) (*models.WeatherApi, error) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&key=%s",
		city, c.weatherApiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request weather for %s: %w", city, err)
	}
	defer resp.Body.Close()

	var data models.WeatherApi
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode weather response for %s: %w", city, err)
	}

	return &data, nil
}
