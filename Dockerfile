# Distributed Cache - Dockerfile

# Use Ubuntu as base image for C++ compilation
FROM ubuntu:20.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

# Install required packages
RUN apt-get update && apt-get install -y \
    g++ \
    make \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy source code
COPY src/ ./src/
COPY client/ ./client/
COPY Makefile ./

# Build the application
RUN make

# Expose port 8080 for the server
EXPOSE 8080

# Default command (can be overridden)
CMD ["./bin/server"]