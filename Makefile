ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

include tools/tools.mk

.PHONY: check
check: lint test

.PHONY: fmt
fmt: license-header-add
	$(GOLANGCI_LINT) run --fix

.PHONY: lint
lint: license-header-check
	$(GOLANGCI_LINT) run

sourcefiles := $(shell find . -type f -name "*.go")

license-header-add:
	@echo "Add License Headers to source files"
	@go run github.com/google/addlicense -c "Nitric Technologies Pty Ltd." -y "2023" $(sourcefiles)

license-header-check:
	@echo "Checking License Headers for source files"
	@go run github.com/google/addlicense -check -c "Nitric Technologies Pty Ltd." -y "2023" $(sourcefiles)

license-check:
	@echo Checking OSS Licenses
	@go build -o ./bin/licenses ./licenses.go 
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/licenses

.PHONY: clean
clean:
	@rm -rf ./interfaces

check-gopath:
ifndef GOPATH
  $(error GOPATH is undefined)
endif

generate:
	go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1 KvStoreClient,KvStore_ScanKeysClient > mocks/keyvalue.go
	go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/proto/queues/v1 QueuesClient > mocks/queues.go
	go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/proto/storage/v1 StorageClient > mocks/storage.go
	go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/proto/secrets/v1 SecretManagerClient > mocks/secrets.go
	go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/proto/topics/v1 TopicsClient > mocks/topics.go
	go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/proto/batch/v1 BatchClient > mocks/batch.go
	go run github.com/golang/mock/mockgen -package mock_v1 google.golang.org/grpc ClientConnInterface > mocks/grpc_clientconn.go
	go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/keyvalue KeyValue,Store > mocks/mockapi/keyvalue.go
	go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/queues Queues,Queue > mocks/mockapi/queues.go
	go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/secrets Secrets,SecretRef > mocks/mockapi/secrets.go
	go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/storage Storage,Bucket > mocks/mockapi/storage.go
	go run github.com/golang/mock/mockgen -package mockapi github.com/nitrictech/go-sdk/api/batch Batch,Job > mocks/mockapi/batch.go

# Runs tests for coverage upload to codecov.io
test-ci: generate
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./resources/... ./api/... ./faas/...

.PHONY: test
test: generate
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover ./resources/... ./api/... ./faas/...