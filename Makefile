export GO111MODULE=on

.PHONY: up
up:
	docker-compose up -d pg mssql bouncer

.PHONY: down
down:
	docker-compose down

.PHONY: test
test:
	go test -v ./...
	cd test && go test -v ./...

example: up
	cd example && go run .
