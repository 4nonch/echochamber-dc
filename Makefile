ifeq ($(wildcard .env), .env)
	include .env
	export
endif

.PHONY: run tests

run:
	go run ./src/*.go
tests:
	go test ./src/...
