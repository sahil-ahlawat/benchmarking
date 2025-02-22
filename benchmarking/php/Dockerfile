# Use PHP 8.2-FPM base image
FROM php:8.2-fpm  

# Install required packages
RUN apt-get update && apt-get install -y nginx supervisor wget

# Set working directory
WORKDIR /var/www/html

# Copy project files
COPY . /var/www/html

# Set permissions
RUN chown -R www-data:www-data /var/www/html

# Copy Nginx config
COPY nginx.conf /etc/nginx/sites-available/default

# Ensure symbolic link exists for Nginx site
RUN ln -sf /etc/nginx/sites-available/default /etc/nginx/sites-enabled/default

# Copy Supervisor config
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Install Nginx Prometheus Exporter
RUN wget https://github.com/nginxinc/nginx-prometheus-exporter/releases/latest/download/nginx-prometheus-exporter-linux-amd64 \
    && chmod +x nginx-prometheus-exporter-linux-amd64 \
    && mv nginx-prometheus-exporter-linux-amd64 /usr/local/bin/nginx-prometheus-exporter

# Expose Nginx (port 80) and Prometheus Exporter (port 9113)
EXPOSE 80 9113

# Start Supervisor to manage Nginx, PHP-FPM, and Nginx Exporter
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
