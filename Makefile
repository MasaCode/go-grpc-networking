.PHONY: build install sh run down


build:
	docker-compose -f .docker/docker-compose.yml build

up:
	docker-compose -f .docker/docker-compose.yml up -d

down:
	docker-compose -f .docker/docker-compose.yml down

exec:
	docker-compose -f .docker/docker-compose.yml exec go-grpc bash

