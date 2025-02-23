package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define the metrics
var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nginx_total_requests",
			Help: "Total number of HTTP requests",
		},
		[]string{"status_code"}, // This will differentiate between 2xx, 4xx, 5xx etc.
	)

	up = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nginx_up",
			Help: "Whether the Nginx service is up (1 for up, 0 for down)",
		},
	)
)

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(up)
}

func main() {
	// Set Nginx status page URL
	nginxStatusURL := "http://localhost/nginx_status"

	// Set up HTTP server for exposing Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	// Check Nginx status page periodically
	go func() {
		for {
			// Fetch Nginx status data
			err := fetchNginxStatus(nginxStatusURL)
			if err != nil {
				log.Printf("Error fetching Nginx status: %v", err)
				up.Set(0) // Nginx is down
			} else {
				up.Set(1) // Nginx is up
			}

			// Sleep for 15 seconds before checking again
			time.Sleep(15 * time.Second)
		}
	}()

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// fetchNginxStatus fetches Nginx status page and updates metrics
func fetchNginxStatus(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch Nginx status: %v", err)
	}
	defer resp.Body.Close()

	// Example of response format (text-based):
	// Active connections: 291 
	// server accepts handled requests
	// 111111 111111 222222
	// Reading: 0 Writing: 1 Waiting: 2
	var total, success, error4xx, error5xx int
	_, err = fmt.Fscanf(resp.Body, "Active connections: %d\n", &total)
	if err != nil {
		return fmt.Errorf("failed to parse Nginx status: %v", err)
	}

	// Dummy data for requests
	// In a real scenario, extract real counts from the status page response
	success = 1000
	error4xx = 50
	error5xx = 20

	// Update Prometheus metrics
	totalRequests.WithLabelValues("2xx").Add(float64(success))
	totalRequests.WithLabelValues("4xx").Add(float64(error4xx))
	totalRequests.WithLabelValues("5xx").Add(float64(error5xx))

	log.Printf("Updated metrics: %d total requests, %d success (2xx), %d 4xx errors, %d 5xx errors", total, success, error4xx, error5xx)

	return nil
}
