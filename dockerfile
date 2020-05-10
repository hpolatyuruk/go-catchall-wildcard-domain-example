# Start from golang base image
FROM golang:latest as builder

COPY . /go/src/go-catchall-wildcard-domain-example/app
WORKDIR /go/src/go-catchall-wildcard-domain-example/app 

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /go/src/go-catchall-wildcard-domain-example/app/main .

# Expose port to the outside world
EXPOSE 80
ENTRYPOINT ["./main"] --port 80