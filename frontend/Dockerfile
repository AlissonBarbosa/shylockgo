# Use the Nginx image from Docker Hub
FROM nginx:latest

# Copy your HTML file into the Nginx container
COPY ./index.html /usr/share/nginx/html/index.html

COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80
