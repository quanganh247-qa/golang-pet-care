# Start with a full Golang image to build the application
FROM registry.semaphoreci.com/golang:1.21 as builder

ENV APP_HOME /go/src/bpm
WORKDIR "$APP_HOME"
COPY ./ .

# Build the application
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bpm .

# Use a minimal base image for the final stage
FROM alpine:latest  
# FROM registry.semaphoreci.com/golang:1.21

ENV APP_HOME /go/src/run
ENV APP_BUILDER /go/src/bpm
RUN mkdir -p "$APP_HOME"
# Set the working directory and run the application
# Copy the binary from the builder stage
COPY --from=builder "$APP_BUILDER"/bpm "$APP_HOME"/bpm

RUN mkdir -p "$APP_HOME"/app/static
# Copy static files and environment files if needed
COPY --from=builder  "$APP_BUILDER"/app/static "$APP_HOME"/app/static
COPY  --from=builder "$APP_BUILDER"/app.env "$APP_HOME"/app.env

WORKDIR "$APP_HOME"


EXPOSE 8088
CMD ["./bpm"]