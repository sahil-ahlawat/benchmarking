FROM php:8.2-fpm  
# Ensure youre using a PHP-FPM base image

# Install required packages
RUN apt-get update && apt-get install -y nginx supervisor

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

# Expose ports
EXPOSE 80

# Start Supervisor to manage both Nginx and PHP-FPM
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
