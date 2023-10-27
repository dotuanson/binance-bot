build:
	go build -o cmd/binance-bot/binance-bot cmd/binance-bot/main.go

clean:
	rm -rf ./cmd/binance-bot/binance-bot

deploy: clean build
	./cmd/binance-bot/binance-bot

test:
	go test -v -cover ./test/

.PHONY: build clean deploy build_linux test