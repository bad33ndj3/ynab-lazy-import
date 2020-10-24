GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=ynab-lazy-import
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
fmt:
		@gofumpt -w --extra .
test:
		$(GOTEST) -v -coverprofile=reports/coverage.out ./...
coverage:
		$(GOCMD) tool cover -html=reports/coverage.out
clean:
		$(GOCMD) clean
run:
		$(GOCMD) run main.go api