package handlers

import (
	"encoding/json"
	"go-first-project/services"
	"net/http"
)

type WeatherHandler struct {
	WeatcherService services.WeatcherService
}

func (h WeatherHandler) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var city string = r.URL.Query().Get("city")

	data, err := h.WeatcherService.GetWeather(ctx, city)

	if err != nil {
		http.Error(w, "City not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
