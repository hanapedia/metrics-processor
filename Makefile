# Variables
DOCKER_HUB_USERNAME := hiroki11hanada
REPOSITORY_NAME := metrics-processor
TAG ?= latest
IMAGE_NAME := $(DOCKER_HUB_USERNAME)/$(REPOSITORY_NAME):$(TAG)

# Conditionally set the BUILD_COMMAND and PUSH_COMMAND
BUILD_COMMAND=docker buildx build --platform linux/amd64 -t $(IMAGE_NAME) . --load
BUILD_COMMAND_ARM64=docker buildx build --platform linux/arm64 -t $(IMAGE_NAME) . --load
PUSH_COMMAND=docker push $(IMAGE_NAME)

.PHONY: dev prod push format
dev:
	$(BUILD_COMMAND_ARM64)

prod:
	$(BUILD_COMMAND)
	$(PUSH_COMMAND)

push:
	$(PUSH_COMMAND)

format:
	goimports -w -l $$(find . -type f -name '*.go' -not -path "./vendor/*")
