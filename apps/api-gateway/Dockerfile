# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install swag for generating docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate swagger docs
RUN swag init -g main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway ./main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway .

# Expose port
EXPOSE 8000

# Command to run
CMD ["./api-gateway"]