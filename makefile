ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: default
default: lint test

.PHONY: install
install:
	@echo installing go dependencies
	@go mod download

.PHONY: install-tools
install-tools: install
	@echo Installing tools from tools.go
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get %

.PHONY: init
init: install-tools
	@echo Installing git hooks
	@find .git/hooks -type l -exec rm {} \; && find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

.PHONY: fmt
fmt: license-header-add
	$(GOLANGCI_LINT) run --fix

.PHONY: lint
lint: license-header-check
	$(GOLANGCI_LINT) run

sourcefiles := $(shell find . -type f -name "*.go")

license-header-add:
	@echo "Add License Headers to source files"
	@go run github.com/google/addlicense -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)

license-header-check:
	@echo "Checking License Headers for source files"
	@go run github.com/google/addlicense -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)

license-check: install-tools
	@echo Checking OSS Licenses
	@go build -o ./bin/licenses ./licenses.go 
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/licenses

.PHONY: clean
clean:
	@rm -rf ./interfaces

# Generate mock implementations
generate-mocks:
	@echo Generating Mock RPC Clients
	@go run github.com/golang/mock/mockgen github.com/nitrictech/apis/go/nitric/v1 DocumentServiceClient,EventServiceClient,TopicServiceClient,QueueServiceClient,StorageServiceClient,FaasServiceClient,FaasService_TriggerStreamClient,DocumentService_QueryStreamClient,SecretServiceClient > mocks/clients.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/apis/go/nitric/v1 DocumentServiceServer,EventServiceServer,TopicServiceServer,QueueServiceServer,StorageServiceServer,FaasServiceServer,FaasService_TriggerStreamServer,DocumentService_QueryStreamServer,SecretServiceServer > mocks/servers.go

# Runs tests for coverage upload to codecov.io
test-ci: generate-mocks
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./api/... ./faas/...

test-examples: generate-mocks
	@echo Testing Nitric Go SDK Examples
	@go test -timeout 30s ./examples/...

.PHONY: test
test: generate-mocks
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover ./api/... ./faas/...