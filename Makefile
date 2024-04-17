GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=bin/srs-sip
MAIN_PATH=main/main.go

default: build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME)

install:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	mv $(BINARY_NAME) /usr/local/bin

.PHONY: clean
