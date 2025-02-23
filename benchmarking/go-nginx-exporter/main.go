package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Prometheus metrics
var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nginx_total_requests",
			Help: "Total number of HTTP requests by status code",
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
	nginxLogPath := "/var/log/nginx/access.log"

	// Start Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start Nginx log monitoring
	go func() {
		for {
			err := parseNginxLogs(nginxLogPath)
			if err != nil {
				log.Printf("Error parsing Nginx logs: %v", err)
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

// parseNginxLogs extracts request counts from access.log
func parseNginxLogs(logPath string) error {
	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	statusCounts := make(map[string]int)

	// Regular expression to extract HTTP status codes from logs
	statusRegex := regexp.MustCompile(`\s(\d{3})\s`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := statusRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			statusCode := matches[1]
			statusCounts[statusCode]++
		}
	}

	// Update Prometheus metrics
	for status, count := range statusCounts {
		totalRequests.WithLabelValues(status).Add(float64(count))
	}

	log.Printf("Status Counts: %v", statusCounts)
	return nil
}
