FROM golang:1.23.10-alpine3.22 AS build

WORKDIR /src

# Copy all the source code
COPY ../../ /src
ENV GOPATH=/root/go
ENV PATH="$PATH:/root/go/bin"

# Install build dependencies
RUN apk add --no-cache \
  git \
  go=~1.24 \
  make \
  musl-dev

# Build the app
RUN make install-tools && \
  make install-tailwindcss && \
  make build && \
  mkdir /app && \
  mv /src/tmp/* /app && \
  mv /src/urchin_config.toml /app/ && \
  mv /src/images /app/ && \
  mv /src/static /app/ && \
  mv /src/migrations /app/migrations

FROM alpine:3.22.0

COPY --from=build /app /app
COPY --from=build /root/go/bin /root/go/bin
ENV PATH="$PATH:/root/go/bin:/usr/local/bin:/usr/bin:/sbin"
