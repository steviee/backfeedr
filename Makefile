.PHONY: build install test clean

BINARY_NAME=backfeedr
BUILD_DIR=./cmd/backfeedr
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME) $(BUILD_DIR)

install: build
	cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)

dev:
	go run $(BUILD_DIR)

docker:
	docker build -t $(BINARY_NAME):latest .

run: build
	./$(BINARY_NAME)
