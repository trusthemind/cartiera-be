# Use the Go 1.20 Alpine image as the base image
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /usr/local/src

# Install necessary packages for building the Go application
RUN apk --no-cache add bash make gcc gettext build-base

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./
COPY .env ./
RUN go mod download

# Copy the entire project into the container
COPY . ./

# Build the Go application
RUN go build -o ./bin/app ./main.go || (echo "Build failed" && exit 1)

# Prepare the final stage for running the application
FROM alpine as runner

# Copy the compiled application from the builder stage
COPY --from=builder /usr/local/src/bin/app /app

# Optionally copy configuration files
COPY ./docker-compose.yaml /config.yaml
COPY ./.env /.env

# Set the entry point for the container to run the Go application
ENTRYPOINT ["/app"]

# Optionally specify default arguments for the application
CMD []
