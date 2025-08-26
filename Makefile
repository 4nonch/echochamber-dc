ifeq ($(wildcard .env), .env)
	include .env
	export
endif

.PHONY: run

run:
	go run ./src/*.go
