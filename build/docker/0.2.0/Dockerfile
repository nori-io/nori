# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13.4 as builder
# Add Maintainer Info
LABEL maintainer="nori <info@nori.io>"
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build Nori
RUN go build -o bin/nori  ./cmd/nori.go


# Create ubuntu 18.04 image
FROM ubuntu:18.04 

RUN apt-get install
WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/nori .
RUN chmod +x nori

# Copy defualt config file
RUN mkdir /app/config
COPY /configs/docker/nori.yml /app/config/nori.yml
COPY /build/docker/0.2.0/files/nori.sh /app/nori.sh
RUN chmod +x nori.sh

# Pass Nori Config flag
ARG NORI_CONFIG
ENV NORI_CONFIG=${NORI_CONFIG:-'file;/app/config/nori.yml'}

# Expose port 80, 8080, 29876 to the outside world
EXPOSE 80 8080 29876

# Declare volumes to mount
VOLUME /plugins
VOLUME /configs

# Command to run the executable
ENTRYPOINT ["/app/nori.sh"]