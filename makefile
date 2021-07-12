.PHONY: build
build:
	go build -v ./cmd/web/main.go

.DEFAULT_GOAL := build