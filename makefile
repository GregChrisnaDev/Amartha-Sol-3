export LINTER_VERSION ?= 1.55.0

GO_PACKAGES ?= $(shell go list ./... | grep -v 'examples\|qtest\|mock')
CMDS 		= $(shell cd cmd && find * -name 'main.go' -print)
ODIR        := deploy/_output
UNAME       := $(shell uname)
CUR_DIR 	= $(shell pwd)

# go build environment
CGO_ENABLED := 1
GOOS		:= linux
GOARCH		:= amd64

export VAR_SERVICES ?= $(foreach path, $(CMDS), $(path:%/main.go=%))

ENV_FILE = .env
ifneq ($(wildcard $(ENV_FILE)),)
include .env
export $(shell sed 's/=.*//' .env)
endif

DB_DEFAULT=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?${DB_QUERY_STRING}
DB := $(if $(url),$(url),$(DB_DEFAULT))

# mac osx build
ifeq ($(UNAME), Darwin)
	GOOS 		:= darwin
	GOARCH		:= arm64
endif

MIGRATION_TOOL_EXISTS = 0
ifneq ("$(wildcard $(CUR_DIR)/bin/migrate)","")
    MIGRATION_TOOL_EXISTS = 1
endif

bin:
	@mkdir -p bin

tool-lint: bin
	@test -e ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v${LINTER_VERSION}

tool-migrate: bin
ifeq ($(MIGRATION_TOOL_EXISTS), 1)
	@echo "Migration tool has been existed."
else ifeq ($(UNAME), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else ifeq ($(UNAME), Darwin)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else
	@echo "Your OS is not supported."
endif

lint: tool-lint
	./bin/golangci-lint run -v --timeout 3m0s

test:
	@go test -race -v ${GO_PACKAGES}

migration: tool-migrate
	${CUR_DIR}/bin/migrate create -ext sql -dir ${CUR_DIR}/migrations/ $(name)

migrate: tool-migrate
	${CUR_DIR}/bin/migrate -path ${CUR_DIR}/migrations/ -database "$(DB)" -verbose up

run:
	go run cmd/main.go
