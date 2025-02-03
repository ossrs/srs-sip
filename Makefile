GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=bin/srs-sip
MAIN_PATH=main/main.go
VUE_DIR=html/NextGB

default: build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(VUE_DIR)/dist
	rm -rf $(VUE_DIR)/node_modules

run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME)

install:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	mv $(BINARY_NAME) /usr/local/bin

vue-install:
	cd $(VUE_DIR) && npm install

vue-build:
	cd $(VUE_DIR) && npm run build

vue-dev:
	cd $(VUE_DIR) && npm run dev

all: build vue-build

.PHONY: clean vue-install vue-build vue-dev all
