#build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

LABEL description="Pet Care Service Backend"
LABEL maintainer="Dinh Huu Quang Anh"
LABEL version="1.0.0"

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and config
COPY --from=builder /app/main .
COPY --from=builder /app/app/db/migration ./migration
COPY app.env .

# Expose port
EXPOSE 8088

# Command to run
CMD ["./main"]
