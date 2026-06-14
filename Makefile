IMAGE_NAME=todo-scheduler
IMAGE_TAG=v1
PORT=7540

.PHONY: run dev build docker-build docker-run test

run:
	go run ./cmd

dev:
	CompileDaemon -directory="./cmd" -command="./cmd/cmd"

build:
	CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run:
	docker run -p $(PORT):$(PORT) $(IMAGE_NAME):$(IMAGE_TAG)

test:
	go test ./...
