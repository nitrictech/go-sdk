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
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

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
	@go run github.com/golang/mock/mockgen github.com/nitrictech/apis/go/nitric/v1 DocumentServiceClient,EventServiceClient,TopicServiceClient,QueueServiceClient,StorageServiceClient,FaasServiceClient,FaasService_TriggerStreamClient,DocumentService_QueryStreamClient,SecretServiceClient,ResourceServiceClient > mocks/clients.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/apis/go/nitric/v1 DocumentServiceServer,EventServiceServer,TopicServiceServer,QueueServiceServer,StorageServiceServer,FaasServiceServer,FaasService_TriggerStreamServer,DocumentService_QueryStreamServer,SecretServiceServer > mocks/servers.go
	@go run github.com/golang/mock/mockgen -package mock_v1 google.golang.org/grpc ClientConnInterface > mocks/grpc_clientconn.go
	@go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/storage Storage,Bucket > mocks/mockapi/storage.go
	@go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/documents Documents,CollectionRef > mocks/mockapi/documents.go
	@go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/queues Queues,Queue > mocks/mockapi/queues.go
	@go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/events Events,Topic > mocks/mockapi/events.go
	@go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/secrets Secrets,SecretRef > mocks/mockapi/secrets.go

# Runs tests for coverage upload to codecov.io
test-ci: generate-mocks
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./resources/... ./api/... ./faas/...

test-examples: generate-mocks
	@echo Testing Nitric Go SDK Examples
	@go test -timeout 30s ./examples/...

.PHONY: test
test: generate-mocks
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover ./resources/... ./api/... ./faas/...