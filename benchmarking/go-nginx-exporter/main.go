package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	successfulRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "nginx_successful_requests",
			Help: "Total number of successful HTTP requests",
		},
	)
	errorRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "nginx_error_requests",
			Help: "Total number of error HTTP requests",
		},
	)
	up = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nginx_up",
			Help: "Whether the Nginx service is up (1 for up, 0 for down)",
		},
	)
)

var prevHandled, prevTotal int

// Register Prometheus metrics
func init() {
	prometheus.MustRegister(successfulRequests)
	prometheus.MustRegister(errorRequests)
	prometheus.MustRegister(up)
}

func main() {
	nginxStatusURL := "http://localhost:8080/nginx_status"

	// Start Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start Nginx status monitoring
	go func() {
		for {
			err := fetchNginxStatus(nginxStatusURL)
			if err != nil {
				log.Printf("Error fetching Nginx status: %v", err)
				up.Set(0)
			} else {
				up.Set(1)
			}
			time.Sleep(10 * time.Second) // Adjust interval as needed
		}
	}()

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":9114", nil))
}

// fetchNginxStatus extracts request counts from the nginx_status page
func fetchNginxStatus(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch Nginx status: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Regex to extract accepts, handled, total requests
	statusRegex := regexp.MustCompile(`\d+`)
	matches := statusRegex.FindAllString(string(body), -1)
	if len(matches) < 3 {
		return fmt.Errorf("unexpected nginx_status format")
	}

	handled, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("failed to parse handled requests: %v", err)
	}

	total, err := strconv.Atoi(matches[2])
	if err != nil {
		return fmt.Errorf("failed to parse total requests: %v", err)
	}

	// Only increment if new values are greater than previous values
	if handled > prevHandled {
		successfulRequests.Add(float64(handled - prevHandled))
	}
	if total > prevTotal {
		errorRequests.Add(float64((total - handled) - (prevTotal - prevHandled)))
	}

	prevHandled = handled
	prevTotal = total

	log.Printf("Total Requests: %d, Successful: %d, Errors: %d", total, handled, total-handled)
	return nil
}
