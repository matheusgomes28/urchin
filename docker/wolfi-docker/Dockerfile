FROM cgr.dev/chainguard/wolfi-base

ARG go_version=1.21
WORKDIR /app

# Install the necessary dependencies to build things
RUN apk add go-fips-${go_version}

# Install the golang dependencies and get rig of cache
RUN GOBIN=/usr/local/bin go install github.com/pressly/goose/v3/cmd/goose@v3.18.0 \
    && GOBIN=/usr/local/bin go install github.com/a-h/templ/cmd/templ@v0.2.543 \
    && GOBIN=/usr/local/bin go install github.com/cosmtrek/air@v1.49.0 \
    && go clean --cache \
    && go clean --modcache \
    && go clean --testcache \
    && go clean --fuzzcache

USER root

ENTRYPOINT [ "bash" ]
