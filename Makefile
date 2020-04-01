GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=kga
BINARY_MACOS=$(BINARY_NAME)_macos
BINARY_LINUX=$(BINARY_NAME)_linux_amd64

BUILD_DIR=./bin
MAIN_PACKAGE="./cmd/kga"

default: build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PACKAGE)

.PHONY: build-macos
build-macos: fmt
	GOOS="darwin" GOARCH="amd64" $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_MACOS) -v $(MAIN_PACKAGE)

.PHONY: build-linux
build-linux: fmt
	GOOS="linux" GOARCH="amd64" $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_LINUX) -v $(MAIN_PACKAGE)

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -drf $(BUILD_DIR)
