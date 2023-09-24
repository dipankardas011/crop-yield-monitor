# Use the official MySQL image as the base image
FROM mysql:8.1.0

# Environment variables for MySQL configuration
ENV MYSQL_ROOT_PASSWORD=1234
# ENV MYSQL_DATABASE=auth

# Copy the init.sql script into the container
COPY mysql.cnf /etc/my.cnf

COPY init.sql /docker-entrypoint-initdb.d/

# Expose the MySQL port
EXPOSE 3306

# The Docker container will automatically start the MySQL server
