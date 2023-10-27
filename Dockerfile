FROM golang:1.20
WORKDIR /binance-bot
COPY . .
RUN go build -o ./binance-bot cmd/binance-bot/main.go
RUN useradd -u 1000 binance
USER myuser
CMD ["./binance-bot"]
