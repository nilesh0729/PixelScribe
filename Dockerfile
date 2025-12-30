# Build Stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api
# Download migrate tool
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./db/migration
COPY start.sh .
COPY app.env.example .env

# Ensure start.sh is executable
RUN chmod +x start.sh

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]
