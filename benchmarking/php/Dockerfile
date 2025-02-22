# Use PHP 8.2-FPM with Nginx
FROM php:8.2-fpm

# Install Nginx and required tools
RUN apt-get update && apt-get install -y nginx wget curl

# Set working directory
WORKDIR /var/www/html

# Copy project files
COPY . /var/www/html

# Set correct permissions
RUN chown -R www-data:www-data /var/www/html

# Copy Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Install & Setup Nginx Prometheus Exporter
# Install & Setup Nginx Prometheus Exporter (Using Provided tar.gz)
RUN wget -O /tmp/nginx-prometheus-exporter.tar.gz https://github.com/nginxinc/nginx-prometheus-exporter/releases/latest/download/nginx-prometheus-exporter_1.4.1_linux_amd64.tar.gz \
    && tar -xvzf /tmp/nginx-prometheus-exporter.tar.gz -C /tmp \
    && mv /tmp/nginx-prometheus-exporter /usr/local/bin/nginx-prometheus-exporter \
    && chmod +x /usr/local/bin/nginx-prometheus-exporter

# Expose required ports
EXPOSE 80 9113

# Start Nginx, PHP-FPM, and Nginx Exporter together
CMD service nginx start && php-fpm -F & /usr/local/bin/nginx-prometheus-exporter -nginx.scrape-uri=http://localhost/nginx_status
