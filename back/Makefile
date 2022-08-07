VERSION     ?= v1.0.1
BUILD_DATE  = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: build

build:
	go build -ldflags "-w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}"
