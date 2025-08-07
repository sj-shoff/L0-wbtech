.PHONY: build run migrate send-order test-api migrate-down

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

docker-clear:
	docker-compose down -v --rmi all
	docker system prune -a -f --volumes
	sudo systemctl restart docker

docker-rebuild:
	docker-compose down -v
	docker-compose build
	docker-compose up