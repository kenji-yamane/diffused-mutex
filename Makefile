.PHONY: build build-cs build-all p1 p2 p3

CMD:=./cmd
BIN:=./bin

PROCESS:=process
SHARED_RESOURCE:=sharedresource

build:
	go build -o $(BIN)/$(PROCESS) $(CMD)/$(PROCESS).go

build-cs:
	go build -o $(BIN)/$(SHARED_RESOURCE) $(CMD)/$(SHARED_RESOURCE).go

build-all: build build-cs

THREE_SIMULATOR=10004 10003 10002

p1:
	$(BIN)/$(PROCESS) 1 $(THREE_SIMULATOR)

p2:
	$(BIN)/$(PROCESS) 2 $(THREE_SIMULATOR)

p3:
	$(BIN)/$(PROCESS) 3 $(THREE_SIMULATOR)

cs:
	$(BIN)/$(SHARED_RESOURCE)
