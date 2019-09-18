#version:1.0
all: build

build:
	go build -o olivine main.go

rundev:
	go run main.go