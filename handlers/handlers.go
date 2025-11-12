package handlers

import (
	"context"
	"go-first-project/services"
	"net/http"
	"time"
	"encoding/json"


)


type WeatherHandler struct {
	WeatcherService services.WeatcherService
}

func (h WeatherHandler) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	client := &http.Client{Timeout: 5 * time.Second}

	var city string = r.URL.Query().Get("city")

	data, err := h.WeatcherService.GetWeather(ctx, client, city)

	if err != nil {
			http.Error(w, "City not found", http.StatusNotFound)
			return
	}

	json.NewEncoder(w).Encode(data)

}