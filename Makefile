.PHONY: build run migrate send-order test-api

build:
	docker-compose build

run:
	docker-compose up

migrate:
	docker-compose run --rm migrator

send-order:
	tr -d '\n' < ./backend/order.json | docker-compose exec -T kafka sh -c 'kafka-console-producer --broker-list kafka:9092 --topic orders'

test-api:
	curl -s http://localhost:8081/order/b563feb7b2b84b6test