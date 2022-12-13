BIN_CND := "./bin/previewer"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CND) -ldflags "$(LDFLAGS)" ./cmd/previewer  

run: build
	$(BIN_CND) -config ./configs/config.toml
	
up:
	docker-compose -f ./deployments/docker-compose.yaml up

down:
	docker-compose -f ./deployments/docker-compose.yaml down

version: build
	$(BIN_CND) version

test:
	go test -race -count 100 ./...

integration-tests:
	set -e ;\
	docker-compose -f ./deployments/docker-compose.tests.yaml  up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f ./deployments/docker-compose.tests.yaml  run integration_tests go test -race -v --tags=integration ./... || test_status_code=$$? ;\
	docker-compose -f ./deployments/docker-compose.tests.yaml  down ;\
	exit $$test_status_code ;

integration-tests-cleanup:
	docker-compose -f ./deployments/docker-compose.tests.yaml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
	cd deployments ; \
  	docker-compose rm -f	

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
