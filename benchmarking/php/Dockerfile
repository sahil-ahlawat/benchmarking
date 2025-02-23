# Use PHP 8.2-FPM with Nginx
FROM php:8.2-fpm

# Install Nginx and required tools
RUN apt-get update && apt-get install -y nginx wget curl \
    && rm -rf /var/lib/apt/lists/*

# Install Go to build the nginx_exporter_go binary (optional, if needed)
RUN apt-get install -y golang

# Set working directory
WORKDIR /var/www/html

# Copy project files
COPY . /var/www/html

# Set correct permissions
RUN chown -R www-data:www-data /var/www/html
RUN chmod -R 755 /var/www/html

# Copy Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf

# Copy the nginx_exporter_go binary
COPY nginx_exporter_go /usr/local/bin/nginx_exporter_go

# Expose required ports
EXPOSE 80 8081

# Start Nginx, PHP-FPM, and the Nginx Exporter Go
CMD service nginx start && php-fpm -F & /usr/local/bin/nginx_exporter_go
