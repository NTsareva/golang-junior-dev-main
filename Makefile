DOCKER_COMPOSE = ./docker-compose.yml
IMAGE_NAME = golang-junior-dev
IMAGE_TAG = latest

build:
	go build ./cmd/exchanges

test:
	go test -v -cover ./...

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-up:
	docker-compose -f $(DOCKER_COMPOSE) up

docker-down:
	docker-compose -f $(DOCKER_COMPOSE) down

clean:
	go clean
	rm -f exchanges

.PHONY: build test docker-up docker-down clean