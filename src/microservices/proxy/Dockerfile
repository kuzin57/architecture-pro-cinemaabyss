FROM kong:latest

# Set working directory
WORKDIR /usr/local/kong

# Install gettext for envsubst (used for template substitution)
USER root
RUN apt-get update && \
    apt-get install -y --no-install-recommends gettext-base && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy startup script and Kong configuration
COPY --chmod=755 startup.sh /usr/local/bin/startup.sh
COPY kong.yml /usr/local/kong/kong.yml
COPY --chmod=755 docker-entrypoint-wrapper.sh /docker-entrypoint-wrapper.sh

# Switch back to kong user
USER kong

# Expose Kong ports
# 8000 - Proxy port (HTTP)
# 8443 - Proxy port (HTTPS)
# 8001 - Admin API (HTTP)
# 8444 - Admin API (HTTPS)
EXPOSE 8000 8443 8001 8444

# Set environment variables for Kong
# Using DB-less mode for simplicity (can be changed to use PostgreSQL)
ENV KONG_DATABASE=off
ENV KONG_DECLARATIVE_CONFIG=/usr/local/kong/kong.yml
ENV KONG_PROXY_ACCESS_LOG=/dev/stdout
ENV KONG_ADMIN_ACCESS_LOG=/dev/stdout
ENV KONG_PROXY_ERROR_LOG=/dev/stdout
ENV KONG_ADMIN_ERROR_LOG=/dev/stdout
ENV KONG_LOG_LEVEL=info
ENV KONG_ADMIN_LISTEN=0.0.0.0:8001

# Use wrapper entrypoint that runs startup.sh first, then Kong
ENTRYPOINT ["/docker-entrypoint-wrapper.sh"]
CMD ["kong", "docker-start"]
