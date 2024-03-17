#! /usr/bin/env bash

set -euo pipefail

git config --global --add safe.directory '*'

cd /urchin/migrations
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:root@tcp(mariadb:3306)/urchin" goose up

cd /urchin
air -c ./docker/.air.toml
