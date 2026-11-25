FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# Install YottaDB and Go build prerequisites
RUN apt-get update && apt-get install -y --no-install-recommends \
    file \
    binutils \
    libelf-dev \
    libicu-dev \
    libncurses-dev \
    nano \
    wget \
    curl \
    ca-certificates \
    gcc \
    libc6-dev \
    make \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Install YottaDB with UTF-8 support
RUN mkdir -p /tmp/tmp && cd /tmp/tmp \
    && wget -q https://download.yottadb.com/ydbinstall.sh \
    && chmod +x ydbinstall.sh \
    && ./ydbinstall.sh --utf8 default --installdir /opt/yottadb/current \
    && rm -rf /tmp/tmp

ARG TARGETARCH

# Install Go 1.26.3
RUN ARCH="${TARGETARCH:-$(dpkg --print-architecture)}" \
    && wget -q "https://go.dev/dl/go1.26.3.linux-${ARCH}.tar.gz" -O /tmp/go.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm -f /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:/opt/yottadb/current:${PATH}"
ENV PKG_CONFIG_PATH="/opt/yottadb/current:/usr/local/lib/pkgconfig:/usr/share/pkgconfig"
ENV ydb_dist="/opt/yottadb/current"
ENV ydb_chset="UTF-8"

WORKDIR /app
COPY . .

RUN chmod +x /app/entrypoint.sh \
    && . /opt/yottadb/current/ydb_env_set \
    && cd /app/go && go build -o /app/ydb-server . \
    && go clean -cache -modcache \
    && rm -rf /root/.cache /root/go /tmp/*

EXPOSE 8080 8081

ENTRYPOINT ["/app/entrypoint.sh"]
