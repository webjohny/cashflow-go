# Use a builder image
FROM golang:1.18 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN chmod +x ./cmd/entrypoint.sh

# Set the entrypoint script to be executed.
ENTRYPOINT ["./cmd/entrypoint.sh"]

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller base image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]