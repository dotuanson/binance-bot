test:
	go test -v ./test/

build:
	rm -rf ./cmd/binance-bot
	go build -o cmd/binance-bot cmd/main.go

deploy:
	docker compose up --force-recreate --detach --build binance-bot

.PHONY: test build deploy