# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_FOLDER=bin
BINARY_NAME=sql-unit-test
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_FILE=cmd/web/main.go

all: test build
build: 
	$(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_NAME) -v $(MAIN_FILE)
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_FOLDER)/$(BINARY_NAME)
	rm -f $(BINARY_FOLDER)/$(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_NAME) -v $(MAIN_FILE)
	./$(BINARY_FOLDER)/$(BINARY_NAME)
deploy: clean test build-linux
	rsync -av . deploy@ta.do:/opt/sql-unit-test && ssh 'deploy@ta.do' 'supervisorctl restart sql-unit-test'

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_UNIX) -v $(MAIN_FILE)
