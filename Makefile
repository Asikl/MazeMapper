# Makefile for compiling and running Go programs

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
BINARY_NAME=test

# Targets
build:
    $(GOBUILD) -o $(BINARY_NAME) test.go

run:
    $(GORUN) test.go

clean:
    $(GOCLEAN)
    rm -f $(BINARY_NAME)
