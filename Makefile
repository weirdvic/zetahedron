# Makefile

SHELL := /bin/bash

.PHONY: run test
run:
	source .env
	go run ./cmd/api-server

test:
	go test -v ./internal/helpers
