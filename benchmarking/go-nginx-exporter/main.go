package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define Prometheus metrics
var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nginx_total_requests",
			Help: "Total number of HTTP requests",
		},
		[]string{"status_code"},
	)
	up = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nginx_up",
			Help: "Whether the Nginx service is up (1 for up, 0 for down)",
		},
	)
)

// Track previous request count to calculate increments
var prevRequests int

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(up)
}

func main() {
	nginxStatusURL := "http://localhost:8080/nginx_status"

	// Expose metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Periodically fetch Nginx metrics
	go func() {
		for {
			if err := fetchNginxStatus(nginxStatusURL); err != nil {
				log.Printf("Error fetching Nginx status: %v", err)
				up.Set(0) // Mark Nginx as down
			} else {
				up.Set(1) // Mark Nginx as up
			}
			time.Sleep(15 * time.Second)
		}
	}()

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":9114", nil))
}

func fetchNginxStatus(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch Nginx status: %v", err)
	}
	defer resp.Body.Close()

	var activeConnections, reading, writing, waiting int
	var accepted, handled, requests int

	// Parse Nginx status page
	_, err = fmt.Fscanf(resp.Body, "Active connections: %d\n", &activeConnections)
	if err != nil {
		return fmt.Errorf("failed to parse active connections: %v", err)
	}

	_, err = fmt.Fscanf(resp.Body, "server accepts handled requests\n")
	if err != nil {
		return fmt.Errorf("failed to parse request headers: %v", err)
	}

	_, err = fmt.Fscanf(resp.Body, "%d %d %d\n", &accepted, &handled, &requests)
	if err != nil {
		return fmt.Errorf("failed to parse request counts: %v", err)
	}

	_, err = fmt.Fscanf(resp.Body, "Reading: %d Writing: %d Waiting: %d\n", &reading, &writing, &waiting)
	if err != nil {
		return fmt.Errorf("failed to parse connection states: %v", err)
	}

	// Calculate request difference since last check
	increase := requests - prevRequests
	if increase > 0 {
		totalRequests.WithLabelValues("total").Add(float64(increase))
	}
	prevRequests = requests // Update previous count

	log.Printf("Active: %d, Total Requests: %d (new: %d), Reading: %d, Writing: %d, Waiting: %d",
		activeConnections, requests, increase, reading, writing, waiting)

	return nil
}
