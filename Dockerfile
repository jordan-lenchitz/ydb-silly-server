FROM yottadb/yottadb-debian:latest

USER root

# Install Go and build dependencies
RUN apt-get update && apt-get install -y \
    wget \
    gcc \
    make \
    pkg-config \
    && wget https://go.dev/dl/go1.26.3.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.26.3.linux-amd64.tar.gz \
    && rm go1.26.3.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /app
COPY . .

# Build the Go application
# We need to source ydb_env_set for pkg-config to work during build if it's using it,
# but usually we just need the CGO flags.
RUN . /opt/yottadb/current/ydb_env_set && go build -o ydb-server .

ENTRYPOINT ["/bin/bash", "/app/entrypoint.sh"]
