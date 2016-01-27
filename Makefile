PATH=$(HOME)/Workspace/go
BIN=$(GOPATH)/bin
FILES=$(shell find ./ -type f -name "*.go")
OUTPUT_DIR=./bin

export GOPATH=$(PATH)
export GOBIN=$(BIN)

run:
	@echo $(GOPATH)
	$(shell go run $(FILES))

build:
	@echo "Building binary to" $(OUTPUT_DIR)
	$(shell go build -o $(OUTPUT_DIR)/api *.go)
