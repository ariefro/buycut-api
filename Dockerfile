# Build stage
FROM golang:1.22.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN chmod +x wait-for.sh

RUN go build -o main cmd/server/main.go

# Run stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/wait-for.sh .
COPY --from=builder /app/.env.staging .

EXPOSE 8080
CMD [ "/app/main" ]
