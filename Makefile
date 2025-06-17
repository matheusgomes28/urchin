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

prepare_env:
	cp -r migrations tests/helpers/

build: prepare_env install-tailwindcss
	$(TEMPL) generate
	GIN_MODE=release $(GOCMD) build -ldflags "-s" -v -o $(BUILD_DIR)/$(BINARY_NAME) $(URCHIN_DIR)
	GIN_MODE=release $(GOCMD) build -ldflags "-s" -v -o $(BUILD_DIR)/$(ADMIN_BINARY_NAME) $(URCHIN_ADMIN_DIR)
	./tailwindcss -i ./static/style.css -o ./static/css/style.css --minify

test: prepare_env
	$(GOCMD) test -v ./...

clean:
	$(GOCMD) clean
	rm -rf $(BUILD_DIR)

# TODO: For now we support only the linux version of tailwindcss, has to be updated in the future to support Windows and MacOS as well.
install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	go install github.com/a-h/templ/cmd/templ@v0.3.865
	go install github.com/cosmtrek/air@v1.49.0 

install-tailwindcss:
	if [ ! -f tailwindcss ]; then \
		wget -q https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.16/tailwindcss-linux-x64 \
			&& echo "33f254b54c8754f16efbe2be1de38ca25192630dc36f164595a770d4bbf4d893  tailwindcss-linux-x64" | sha256sum -c \
			&& chmod +x tailwindcss-linux-x64 \
			&& mv tailwindcss-linux-x64 tailwindcss; \
	fi

.PHONY: all build test clean install-tailwindcss
