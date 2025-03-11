#!/usr/bin/env bash

set -euo pipefail

# Ensure the admin app uses port 8081
export URCHIN_HTTP_PORT=8081
export URCHIN_ADMIN=true

# Run the admin app with the config file
exec ./tmp/urchin-admin --config ./docker/urchin_config.toml 