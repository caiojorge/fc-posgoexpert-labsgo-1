FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install swag for swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate swagger docs
RUN swag init -g cmd/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather-api ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/weather-api .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./weather-api"]
