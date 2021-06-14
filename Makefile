IMAGE_NAME=url-shortener
IMAGE_TAG=latest

all: lint build-image

lint:
	golangci-lint -c .golangci.yml -v run

build-image: 
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
