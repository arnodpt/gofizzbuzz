include .env.sh

all: run

run:
	docker-compose up -d --build

stop:
	docker-compose down