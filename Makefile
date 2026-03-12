.PHONY: build build-client install test clean dev docker run

BINARY_NAME=backfeedr
CLIENT_NAME=backfeedr-client
BUILD_DIR=./cmd/backfeedr
CLIENT_DIR=./cmd/backfeedr-client
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME) $(BUILD_DIR)

build-client:
	go build -o $(CLIENT_NAME) $(CLIENT_DIR)

build-all: build build-client

install: build
	cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

cp $(CLIENT_NAME) $(INSTALL_PATH)/$(CLIENT_NAME)

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME) $(CLIENT_NAME)

dev:
	go run $(BUILD_DIR)

dev-client:
	go run $(CLIENT_DIR)

docker:
	docker build -t $(BINARY_NAME):latest .

run: build
	./$(BINARY_NAME)
