# Use an official Go image as a parent image
FROM golang:1.17-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Install GCC (required for CGO)
RUN apk add --no-cache gcc musl-dev

# Copy the local package files to the container's workspace.
ADD . /app

# Build your application
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Use a minimal alpine image for the final build
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]
