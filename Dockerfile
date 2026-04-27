# Build stage
FROM golang:1.26.2-alpine3.23 AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o go-shorten cmd/main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/go-shorten .
# Copy static files/templates
COPY --from=builder /app/web ./web
# Copy database schema
COPY --from=builder /app/db ./db

# Expose the port the app runs on
EXPOSE 8000

# Set default environment variables
ENV DATABASE_URL=db/go-shorten.db
ENV PORT=8000
ENV GIN_MODE=release

# Command to run the executable
CMD ["./go-shorten"]
