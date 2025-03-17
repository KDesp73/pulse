package database

import (
	"encoding/json"
	"net/http"
)

func GetAvgTemperature(w http.ResponseWriter, r *http.Request) {
	var avgTemp float32
	query := "SELECT AVG(temperature) FROM plant_data"
	err := db.QueryRow(query).Scan(&avgTemp)
	if err != nil {
		http.Error(w, "Error fetching average temperature", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float32{"average_temperature": avgTemp})
}

// API handler for min/max temperature
func GetMinMaxTemperature(w http.ResponseWriter, r *http.Request) {
	var minTemp, maxTemp float32
	query := "SELECT MIN(temperature), MAX(temperature) FROM plant_data"
	err := db.QueryRow(query).Scan(&minTemp, &maxTemp)
	if err != nil {
		http.Error(w, "Error fetching min/max temperature", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float32{"min_temperature": minTemp, "max_temperature": maxTemp})
}

// API handler for average soil moisture
func GetAvgSoilMoisture(w http.ResponseWriter, r *http.Request) {
	var avgSoil float32
	query := "SELECT AVG(soil) FROM plant_data"
	err := db.QueryRow(query).Scan(&avgSoil)
	if err != nil {
		http.Error(w, "Error fetching average soil moisture", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float32{"average_soil_moisture": avgSoil})
}

// API handler for the most recent reading
func GetLatestReading(w http.ResponseWriter, r *http.Request) {
	var timestamp string
	var temperature, humidity, soil, light float32
	query := "SELECT timestamp, temperature, humidity, soil, light FROM plant_data ORDER BY timestamp DESC LIMIT 1"
	err := db.QueryRow(query).Scan(&timestamp, &temperature, &humidity, &soil, &light)
	if err != nil {
		http.Error(w, "Error fetching latest reading", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timestamp":   timestamp,
		"temperature": temperature,
		"humidity":    humidity,
		"soil":        soil,
		"light":       light,
	})
}
