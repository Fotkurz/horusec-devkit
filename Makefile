GO ?= go
GOFMT ?= gofmt
GO_FILES ?= $$(find . -name '*.go' | grep -v vendor)
GOLANG_CI_LINT ?= ./bin/golangci-lint
GO_IMPORTS ?= goimports
GO_IMPORTS_LOCAL ?= github.com/ZupIT/horusec-devkit
HORUSEC ?= horusec
COMPOSE_FILE_NAME ?= docker-compose.yaml
DOCKER_COMPOSE ?= docker-compose

fmt:
	$(GOFMT) -w $(GO_FILES)

lint:
    ifeq ($(wildcard $(GOLANG_CI_LINT)), $(GOLANG_CI_LINT))
		$(GOLANG_CI_LINT) run -v --timeout=5m -c .golangci.yml ./...
    else
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
		$(GOLANG_CI_LINT) run -v --timeout=5m -c .golangci.yml ./...
    endif

coverage:
	chmod +x scripts/coverage.sh
	scripts/coverage.sh 99.2 "."

test:
	$(GO) clean -testcache && $(GO) test -v ./... -timeout=2m -parallel=1 -failfast -short

fix-imports:
    ifeq (, $(shell which $(GO_IMPORTS)))
		$(GO) get -u golang.org/x/tools/cmd/goimports
		$(GO_IMPORTS) -local $(GO_IMPORTS_LOCAL) -w $(GO_FILES)
    else
		$(GO_IMPORTS) -local $(GO_IMPORTS_LOCAL) -w $(GO_FILES)
    endif

security:
    ifeq (, $(shell which $(HORUSEC)))
		curl -fsSL https://horusec.io/bin/install.sh | bash
		$(HORUSEC) start -p="./" -e="true"
    else
		$(HORUSEC) start -p="./" -e="true"
    endif

compose: compose-down compose-up

compose-down:
	$(DOCKER_COMPOSE) -f deployments/$(COMPOSE_FILE_NAME) down -v

compose-up:
	$(DOCKER_COMPOSE) -f deployments/$(COMPOSE_FILE_NAME) up -d --build --force-recreate

migrate: migrate-drop migrate-up

migrate-up:
	chmod +x ./scripts/migration-run.sh
	./scripts/migration-run.sh up

migrate-drop:
	chmod +x ./scripts/migration-run.sh
	./scripts/migration-run.sh drop -f

pipeline: fmt fix-imports lint test coverage security
