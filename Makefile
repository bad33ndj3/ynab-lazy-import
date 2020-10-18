GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
BINARY_NAME=csvtoynab
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
fmt:
		$(GOFMT) ./...
test:
		$(GOTEST) -v -coverprofile=reports/coverage.out ./...
coverage:
		$(GOCMD) tool cover -html=reports/coverage.out