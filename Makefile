test:
	go test -v -cover ./test/

build:
	rm -rf ./cmd/binance-bot/binance-bot
	go build -o cmd/binance-bot/binance-bot cmd/binance-bot/main.go

deploy:
	docker compose up

.PHONY: build clean deploy build_linux test