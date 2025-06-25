FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY container.env /app/.env

RUN go build -o ordersystem ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

FROM alpine:latest

COPY --from=builder /app/ordersystem /app/ordersystem
COPY --from=builder /app/.env /app/.env

EXPOSE 8080
EXPOSE 8000
EXPOSE 50051

WORKDIR /app

CMD ["./ordersystem"]
