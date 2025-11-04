FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY src/microservices/events/go.mod ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY src/microservices/events/ .

# Install make
RUN apk add --no-cache make

# Copy the OpenAPI specification
COPY api-specification.yaml ./api-specification.yaml

# Generate code from OpenAPI specification
RUN make generate OPENAPI_SPEC=./api-specification.yaml

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o events-service .

# Use a smaller image for the final application
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/events-service .

# Expose the application port
EXPOSE 8082

# Command to run the executable
CMD ["./events-service"]
