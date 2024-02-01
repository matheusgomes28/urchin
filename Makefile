# Go parameters
GOCMD=go
TEMPL=templ
TAILWIND=tailwindcss-linux-x64
BUILD_DIR=./tmp
URCHIN_DIR=./cmd/urchin
URCHIN_ADMIN_DIR=./cmd/urchin-admin

# Name of the binary
BINARY_NAME=urchin
ADMIN_BINARY_NAME=urchin-admin

all: build test

build:
	$(TEMPL) generate
	$(TAILWIND) -i ./static/style.css -o ./static/output.css -m
	$(GOCMD) build -v -o $(BUILD_DIR)/$(BINARY_NAME) $(URCHIN_DIR)
	$(GOCMD) build -v -o $(BUILD_DIR)/$(ADMIN_BINARY_NAME) $(URCHIN_ADMIN_DIR)

test:
	$(GOCMD) test -v ./...

clean:
	$(GOCMD) clean
	rm -rf $(BUILD_DIR)

.PHONY: all build test clean
