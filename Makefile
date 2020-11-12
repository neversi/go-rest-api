.PHONY: build

build:
	go build -v ./main_program/main.go

.DEFAULT_GOAL := build