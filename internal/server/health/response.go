package health

import "time"

// Response represents the health check response
type Response struct {
	Status    string    `json:"status"`            // Current server status
	Timestamp time.Time `json:"timestamp"`         // Current server time
	Uptime    string    `json:"uptime"`            // Server uptime
	Version   string    `json:"version,omitempty"` // Optional version info
}
