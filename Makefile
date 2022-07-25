include .env.sh

all: run

run:
	docker-compose up -d --build

stop:
	docker-compose down

test:
	go test ./api/http_test.go