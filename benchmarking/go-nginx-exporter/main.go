package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
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

// Register Prometheus metrics
func init() {
	prometheus.MustRegister(totalRequests)
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
			time.Sleep(15 * time.Second) // Adjust interval as needed
		}
	}()

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":9114", nil))
}

// fetchNginxStatus extracts real status counts from the Nginx status page
func fetchNginxStatus(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch Nginx status: %v", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var requests int
	var statusCounts = make(map[string]int)

	// Regular expression to extract status codes
	statusRegex := regexp.MustCompile(`\b(\d{3})\b`)

	// Parse Nginx status page
	lineNumber := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNumber++

		if lineNumber == 3 { // Third line contains request counts
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				requests, err = strconv.Atoi(fields[2])
				if err != nil {
					return fmt.Errorf("failed to parse request count: %v", err)
				}
			}
		}

		// Extract status codes dynamically
		matches := statusRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			statusCode := match[1]
			statusCounts[statusCode]++
		}
	}

	// Update Prometheus metrics
	for status, count := range statusCounts {
		totalRequests.WithLabelValues(status).Add(float64(count))
	}

	log.Printf("Requests: %d | Status Counts: %v", requests, statusCounts)
	return nil
}
