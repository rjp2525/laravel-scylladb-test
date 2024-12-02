FROM golang:latest as build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copy the application code
COPY . .

# Build the Go binary as a statically linked executable
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest as run

WORKDIR /app

# Copy the statically linked Go binary from the build stage
COPY --from=build /app/main /app/main

EXPOSE 8080

# Run the binary
CMD ["./main"]
