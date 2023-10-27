FROM golang:1.20
WORKDIR /binance-bot
COPY . .
RUN go build -o ./binance-bot cmd/main.go
RUN useradd -u 1000 binance-user
USER binance-user
CMD ["./binance-bot"]
