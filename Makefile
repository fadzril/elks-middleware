PATH=$(HOME)/Workspace/go
BIN=$(GOPATH)/bin
FILES=$(shell find ./ -type f -name "*.go")

export GOPATH=$(PATH)
export GOBIN=$(BIN)

run:
	@echo $(GOPATH)
	$(shell go run $(FILES))
