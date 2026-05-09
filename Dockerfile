FROM yottadb/yottadb-debian:latest

USER root

# Install Node.js 20 and dependencies for mg-dbx-napi
RUN apt-get update && apt-get install -y \
    curl \
    g++ \
    make \
    python3 \
    pkg-config \
    libssl-dev \
    procps \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs

WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

# Ensure /data exists and is writable
RUN mkdir -p /data

ENTRYPOINT ["/bin/bash", "/app/entrypoint.sh"]
