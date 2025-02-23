#!/bin/bash

# Start PHP-FPM
php-fpm -D

# Start Nginx
openresty -g "daemon off;" &

# Start Nginx Prometheus Exporter
/usr/local/bin/nginx-prometheus-exporter -nginx.scrape-uri=http://localhost/nginx_status
