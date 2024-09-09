# Use a builder image
FROM golang:1.18 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
#RUN ./.env
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller base image
FROM alpine:latest AS production
WORKDIR /app/
COPY --from=builder /app/main main
COPY ./collections collections
COPY ./professions.json .
EXPOSE 8080
RUN chmod +x ./main  # Ensure binary is executable
RUN ls -l ./
ENTRYPOINT ["./main"]