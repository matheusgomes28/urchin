# Go parameters
GOCMD=go
TEMPL=templ
BUILD_DIR=./tmp
URCHIN_DIR=./cmd/urchin
URCHIN_ADMIN_DIR=./cmd/urchin-admin

# Name of the binary
BINARY_NAME=urchin
ADMIN_BINARY_NAME=urchin-admin

all: build test

build:
	$(TEMPL) generate
	GIN_MODE=release $(GOCMD) build -ldflags "-s" -v -o $(BUILD_DIR)/$(BINARY_NAME) $(URCHIN_DIR)
	GIN_MODE=release $(GOCMD) build -ldflags "-s" -v -o $(BUILD_DIR)/$(ADMIN_BINARY_NAME) $(URCHIN_ADMIN_DIR)

test:
	$(GOCMD) test -v ./...

clean:
	$(GOCMD) clean
	rm -rf $(BUILD_DIR)

install-tools:
	go install github.com/a-h/templ/cmd/templ@v0.2.543

.PHONY: all build test clean
