.PHONY: build install sh run down


up:
	docker-compose -f .docker/docker-compose.yml up -d

down:
	docker-compose -f .docker/docker-compose.yml down

exec:
	docker-compose -f .docker/docker-compose.yml exec go-grpc bash

