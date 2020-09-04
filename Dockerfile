FROM golang:latest as builder
LABEL maintainer="Ron Blom <blom.ron@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
# COPY telegram-sidecar/go.mod telegram-sidecar/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download
RUN go get github.com/go-telegram-bot-api/telegram-bot-api

# Copy the source from the current directory to the Working Directory inside the container
COPY telegram-sidecar/telegram.go .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telegram .


######## Start a new stage from scratch #######
FROM alpine:latest  
RUN apk --no-cache add ca-certificates curl ffmpeg

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/telegram .

# Expose port 8090 to the outside world
EXPOSE 8090

# Command to run the executable
CMD ["./telegram"] 