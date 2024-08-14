#!/bin/bash
set -e

# non-conditional step-down to app user
echo "docker-entrypoint.sh: switching user from root to $APP_USER..."
chown -R "$APP_USER:$APP_GROUP" "$APP_HOME"
exec gosu "$APP_USER" "$@"
