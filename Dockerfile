FROM golang:1.16-alpine as builder

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o /app/server .


FROM alpine:3.13.5

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/server /app/server

# This container exposes port 8080 to the outside world
EXPOSE 80

# Run the server
CMD ["./server"]