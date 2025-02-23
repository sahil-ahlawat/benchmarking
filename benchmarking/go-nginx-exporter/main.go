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
	nginxStatusURL := "http://localhost:8080/nginx_status"

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
	log.Fatal(http.ListenAndServe(":9114", nil))
}

// fetchNginxStatus fetches Nginx status page and updates metrics
func fetchNginxStatus(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch Nginx status: %v", err)
	}
	defer resp.Body.Close()

	// Declare variables to hold parsed values
	var activeConnections, reading, writing, waiting int
	var totalRequestsCount, handledRequests, requests int // Renamed variable

	// Use fmt.Fscanf to extract the data from the Nginx status page
	_, err = fmt.Fscanf(resp.Body, "Active connections: %d\n", &activeConnections)
	if err != nil {
		return fmt.Errorf("failed to parse 'Active connections': %v", err)
	}

	// Skip the next two lines
	_, err = fmt.Fscanf(resp.Body, "server accepts handled requests\n")
	if err != nil {
		return fmt.Errorf("failed to parse 'server accepts handled requests' line: %v", err)
	}
	_, err = fmt.Fscanf(resp.Body, "%d %d %d\n", &totalRequestsCount, &handledRequests, &requests)
	if err != nil {
		return fmt.Errorf("failed to parse request counts: %v", err)
	}

	// Get the last line for Reading, Writing, and Waiting values
	_, err = fmt.Fscanf(resp.Body, "Reading: %d Writing: %d Waiting: %d\n", &reading, &writing, &waiting)
	if err != nil {
		return fmt.Errorf("failed to parse 'Reading Writing Waiting' line: %v", err)
	}

	// Now update your Prometheus metrics with the values
	// Example: Update success/error request counts (you may need to extract actual data)
	success, error4xx, error5xx := 1000, 50, 20 // Replace with actual parsing if available

	// Update Prometheus metrics
	totalRequests.WithLabelValues("2xx").Add(float64(success))
	totalRequests.WithLabelValues("4xx").Add(float64(error4xx))
	totalRequests.WithLabelValues("5xx").Add(float64(error5xx))

	log.Printf("Active connections: %d, Total Requests: %d, Reading: %d, Writing: %d, Waiting: %d", activeConnections, totalRequestsCount, reading, writing, waiting)
	log.Printf("Updated metrics: %d total requests, %d success (2xx), %d 4xx errors, %d 5xx errors", totalRequestsCount, success, error4xx, error5xx)

	return nil
}
