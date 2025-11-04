#!/bin/sh
set -e
sh /usr/local/bin/startup.sh
exec /docker-entrypoint.sh "$@"

