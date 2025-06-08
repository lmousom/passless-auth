package handlers

import (
	"encoding/json"
	"net/http"
)

var (
	version   = "1.0.1"
	buildTime = "2025-06-08"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
	BuildTime string            `json:"build_time"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status: "ok",
		Services: map[string]string{
			"database": "up",
			"cache":    "up",
			"sms":      "up",
		},
		Version:   version,
		BuildTime: buildTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
