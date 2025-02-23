# Use PHP 8.2-FPM with Nginx
FROM php:8.2-fpm

# Install Nginx and required tools
RUN apt-get update && apt-get install -y nginx wget curl \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /var/www/html

# Copy project files
COPY . /var/www/html

# Copy the precompiled nginx_exporter_go binary into the container
COPY nginx_exporter_go /usr/local/bin/nginx_exporter_go

# Set correct permissions
RUN chown -R www-data:www-data /var/www/html
RUN chmod -R 755 /var/www/html
RUN chmod +x /usr/local/bin/nginx_exporter_go

# Copy Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf

# Expose required ports
EXPOSE 80
EXPOSE 9113  # Port for nginx_exporter_go

# Start Nginx and the nginx_exporter_go exporter
CMD service nginx start && /usr/local/bin/nginx_exporter_go & php-fpm -F
