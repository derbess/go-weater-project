package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"go-first-project/services"
	"log"
	"net/http"
)

type WeatherHandler struct {
	WeatcherService services.WeatcherService
}

func (h WeatherHandler) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		log.Printf("Request aborted: server is shutting down")
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	default:
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	data, err := h.WeatcherService.GetWeather(ctx, city)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("Request cancelled for city %s: server shutting down", city)
			http.Error(w, "Service unavailable - shutting down", http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Request timeout for city %s", city)
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			return
		}

		log.Printf("Error getting weather for city %s: %v", city, err)
		http.Error(w, "Failed to get weather data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
