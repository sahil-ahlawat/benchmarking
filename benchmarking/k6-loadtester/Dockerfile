# Use official k6 image
FROM grafana/k6:latest

# Set working directory inside the container
WORKDIR /k6-test

# Copy test script into container
COPY script.js /k6-test/script.js

# Set entrypoint to run k6 script
ENTRYPOINT ["k6", "run"]
