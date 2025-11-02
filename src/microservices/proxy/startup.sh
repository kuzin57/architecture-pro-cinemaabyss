#!/bin/sh

# Startup script for Kong API Gateway
# This script will be executed when the container starts, before Kong starts

# Exit on error
set -e

echo "Running startup script..."

# Заменяем переменные окружения в kong.yml
echo "Substituting environment variables in kong.yml..."

# Проверяем наличие необходимых переменных окружения
if [ -z "$MONOLITH_URL" ]; then
  echo "Warning: MONOLITH_URL is not set, using default: http://monolith:8080"
  export MONOLITH_URL="http://monolith:8080"
fi

if [ -z "$MOVIES_SERVICE_URL" ]; then
  echo "Warning: MOVIES_SERVICE_URL is not set, using default: http://movies-service:8081"
  export MOVIES_SERVICE_URL="http://movies-service:8081"
fi

if [ -z "$EVENTS_SERVICE_URL" ]; then
  echo "Warning: EVENTS_SERVICE_URL is not set, using default: http://events-service:8082"
  export EVENTS_SERVICE_URL="http://events-service:8082"
fi

# Проверяем и устанавливаем значения feature flag
if [ -z "$GRADUAL_MIGRATION" ]; then
  echo "Warning: GRADUAL_MIGRATION is not set, using default: false"
  export GRADUAL_MIGRATION="false"
fi

if [ -z "$MOVIES_MIGRATION_PERCENT" ]; then
  echo "Warning: MOVIES_MIGRATION_PERCENT is not set, using default: 0"
  export MOVIES_MIGRATION_PERCENT="0"
fi

# Преобразуем GRADUAL_MIGRATION в булево значение для Lua
if [ "$GRADUAL_MIGRATION" = "true" ] || [ "$GRADUAL_MIGRATION" = "1" ]; then
  GRADUAL_MIGRATION_LUA="true"
else
  GRADUAL_MIGRATION_LUA="false"
fi

# Используем envsubst для подстановки переменных окружения прямо в kong.yml
# Если envsubst недоступен, используем sed как fallback
if command -v envsubst >/dev/null 2>&1; then
  envsubst < /usr/local/kong/kong.yml > /usr/local/kong/kong.yml.tmp && \
    mv /usr/local/kong/kong.yml.tmp /usr/local/kong/kong.yml
  echo "Substituted variables in kong.yml using envsubst"
else
  # Fallback: используем sed для замены переменных прямо в файле
  sed -i "s|\${MONOLITH_URL}|${MONOLITH_URL}|g; \
          s|\${MOVIES_SERVICE_URL}|${MOVIES_SERVICE_URL}|g; \
          s|\${EVENTS_SERVICE_URL}|${EVENTS_SERVICE_URL}|g; \
          s|\${GRADUAL_MIGRATION}|${GRADUAL_MIGRATION_LUA}|g; \
          s|\${MOVIES_MIGRATION_PERCENT}|${MOVIES_MIGRATION_PERCENT}|g" \
    /usr/local/kong/kong.yml
  echo "Substituted variables in kong.yml using sed"
fi

echo "Kong configuration generated successfully:"
echo "  MONOLITH_URL=${MONOLITH_URL}"
echo "  MOVIES_SERVICE_URL=${MOVIES_SERVICE_URL}"
echo "  EVENTS_SERVICE_URL=${EVENTS_SERVICE_URL}"
echo "  GRADUAL_MIGRATION=${GRADUAL_MIGRATION} (Lua: ${GRADUAL_MIGRATION_LUA})"
echo "  MOVIES_MIGRATION_PERCENT=${MOVIES_MIGRATION_PERCENT}"

echo "Startup script completed. Starting Kong..."

