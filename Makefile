.PHONY: install-ydb test test-mumps bench run mumps-xexe mumps-zsrv docker-build docker-up docker-down

YDB_ENV = YDB_DIST=$$(find /usr/local/lib/yottadb /opt/yottadb -name ydb_env_set -exec dirname {} \; 2>/dev/null | head -n 1); \
	export ydb_chset=UTF-8; \
	export ydb_dir=$${ydb_dir:-$(PWD)/ydb-data}; \
	mkdir -p "$$ydb_dir"; \
	. "$$YDB_DIST/ydb_env_set"; \
	export ydb_routines="$(PWD)/mumps $$ydb_routines"; \
	export PKG_CONFIG_PATH="$$YDB_DIST:/usr/local/lib/pkgconfig:/usr/share/pkgconfig"

install-ydb:
	apt-get update && apt-get install -y --no-install-recommends file binutils libelf-dev libicu-dev libncurses-dev nano wget ca-certificates gcc libc6-dev make pkg-config
	mkdir -p /tmp/tmp && cd /tmp/tmp && wget -q https://download.yottadb.com/ydbinstall.sh && chmod +x ydbinstall.sh && ./ydbinstall.sh --utf8 default --verbose

test-mumps:
	@$(YDB_ENV); "$$ydb_dist/yottadb" -run ^ZTEST

test:
	@$(YDB_ENV); cd go && go test -v -race ./...

bench:
	@$(YDB_ENV); cd go && go test -bench=. -benchmem ./...

run:
	@$(YDB_ENV); cd go && go run .

mumps-xexe:
	@$(YDB_ENV); "$$ydb_dist/yottadb" -run ^XEXE

mumps-zsrv:
	@$(YDB_ENV); "$$ydb_dist/yottadb" -run ^ZSRV

docker-build:
	docker build --platform linux/amd64 -t ydb-silly-server:latest -t ghcr.io/jordan-lenchitz/ydb-silly-server:latest .

docker-push: docker-build
	docker push ghcr.io/jordan-lenchitz/ydb-silly-server:latest

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down
