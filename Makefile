APP_NAME = premo
APP_VERSION = 1.0.0

# build with version infos
VERSION_DIR = github.com/meshplus/${APP_NAME}
BUILD_DATE = $(shell date +%FT%T)
GIT_COMMIT = $(shell git log --pretty=format:'%h' -n 1)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

GO_LDFLAGS += -X "${VERSION_DIR}.BuildDate=${BUILD_DATE}"
GO_LDFLAGS += -X "${VERSION_DIR}.CurrentCommit=${GIT_COMMIT}"
GO_LDFLAGS += -X "${VERSION_DIR}.CurrentBranch=${GIT_BRANCH}"
GO_LDFLAGS += -X "${VERSION_DIR}.CurrentVersion=${APP_VERSION}"

TEST_PKGS := $(shell go list ./... | grep -v 'mock_*' | grep -v 'tester')
TEST_TIME := $(shell date "+%Y%m%d%H%M%S")

RED=\033[0;31m
GREEN=\033[0;32m
BLUE=\033[0;34m
NC=\033[0m

GO = go

help: Makefile
	@printf "${BLUE}Choose a command run:${NC}\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/    /'

## make prepare: Preparation before development
prepare:
	@cd scripts && bash prepare.sh

## make test: Run go unittest
test:
	$(GO) generate ./...
	@$(GO) test ${TEST_PKGS} -count=1

## make test-coverage: Test project with cover
test-coverage:
	$(GO) generate ./...
	@$(GO) test -short -coverprofile cover.out -covermode=atomic ${TEST_PKGS}
	@cat cover.out >> coverage.txt

## make bitxhub-tester: Run bitxhub test
bitxhub-tester:
ifeq ("${REPORT}", "Y")
	cd tester/bxh_tester && $(GO) test -v -run TestTester -json > json-report.txt
	gotest2allure -f json-report.txt
	zip -r allure-results.zip allure-results
else
	cd tester/bxh_tester && $(GO) test -v -run TestTester
endif

## make interchain-tester: Run interchain test
interchain-tester:
	@cd tester/interchain_tester && $(GO) test -v -run TestTester

## make gosdk-tester: Run gosdk test
gosdk-tester:
ifeq ("${REPORT}", "Y")
	cd tester/gosdk_tester && $(GO) test -v -run TestTester -json > json-report.txt
	gotest2allure -f json-report.txt
	zip -r allure-results.zip allure-results
else
	cd tester/gosdk_tester && $(GO) test -v -run TestTester
endif

## make http-tester: Run http test
http-tester:
ifeq ("${REPORT}", "Y")
	cd tester/http_tester && $(GO) test -v -run TestTester -json > json-report.txt
	gotest2allure -f json-report.txt
	zip -r allure-results.zip allure-results
else
	cd tester/http_tester && $(GO) test -v -run TestTester
endif

## make install: Go install the project
install:
	cd internal/repo && packr
	$(GO) install -ldflags '${GO_LDFLAGS}' ./cmd/${APP_NAME}
	@printf "${GREEN}Build ${APP_NAME} successfully!${NC}\n"

build:
	cd internal/repo && packr
	@mkdir -p bin
	$(GO) build -ldflags '${GO_LDFLAGS}' ./cmd/${APP_NAME}
	@mv ./${APP_NAME} bin
	@printf "${GREEN}Build ${APP_NAME} successfully!${NC}\n"

## make linter: Run golanci-lint
linter:
	run

.PHONY: tester build