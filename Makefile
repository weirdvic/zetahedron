# Makefile

SHELL := /bin/bash

.PHONY: run
run:
	source .env
	go run ./cmd/api-server
