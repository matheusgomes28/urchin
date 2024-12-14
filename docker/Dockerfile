# Version: 0.3

#
# We base our image from the alpine light image
FROM ubuntu:jammy-20240125

#
# Environment variables needed for the
# build system
ENV TZ=Europe/London
ENV DEBIAN_FRONTEND="noninteractive"
ENV PATH="${PATH}:/usr/local/go/bin"

#
# Identify the maintainer of an image
LABEL maintainer="matheusgarcia28@gmail.com"

# install build dependencies 
RUN apt-get update \
    && apt-get install -y --no-install-recommends --no-install-suggests \
    ca-certificates \
    git \
    make \
    pkg-config \
    wget \
    && wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz -P golang-install \
    && cd golang-install \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz \
    && GOBIN=/usr/local/bin go install github.com/pressly/goose/v3/cmd/goose@v3.18.0 \
    && GOBIN=/usr/local/bin go install github.com/a-h/templ/cmd/templ@v0.2.543 \
    && GOBIN=/usr/local/bin go install github.com/cosmtrek/air@v1.49.0 \
    && wget https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh -P golangci-install \ 
    && sh ./golangci-install/install.sh -b /usr/local/bin v1.56.2 \
    && rm -rf ./golangci-install \
    && apt-get clean \  
    && apt-get autoremove -y \
    && apt-get autoremove --purge -y \
    && go clean --cache \
    && go clean --modcache \
    && go clean --testcache \
    && go clean --fuzzcache \
    && rm -rf golang-install \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /usr/lib/gcc/x86_64-linux-gnu/7*

