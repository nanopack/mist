DISPLAY_BOLD    := "\033[01m"
DISPLAY_RESET   := "\033[0;0m"

VERSION ?= 1.3

# Where we're running Make from
SRC_DIR ?= $(shell pwd)

# We'll publish all containers out of this repo into deepmap-go scope.
DOCKER_IMAGE_NAME := mist

# This rule builds our Go development environment in a Docker container
.phony: docker

docker:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
