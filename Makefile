build_linux:
	env GOARCH=amd64 GOOS=linux go build -o ./binance-bot cmd/binance-bot/main.go

build:
	go build -o cmd/binance-bot/binance-bot cmd/binance-bot/main.go

clean:
	rm -rf ./cmd/binance-bot/binance-bot

deploy: clean build
	./cmd/binance-bot/binance-bot

.PHONY: build clean deploy build_linux