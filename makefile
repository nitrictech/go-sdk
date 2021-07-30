install:
	@echo installing go dependencies
	@go mod download

install-tools: install
	@echo Installing tools from tools.go
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get %

init: install-tools
	@echo Installing git hooks
	@find .git/hooks -type l -exec rm {} \; && find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

fmt:
	@echo Formatting Code
	@gofmt -s -w ./**/*.go

lint:
	@echo Formatting Code
	@golint ./...

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

clean:
	@rm -rf ./interfaces

# Generate interfaces
generate-proto: clean
	@echo Generating Proto Sources
	@mkdir -p ./interfaces/
	@protoc --go_out=./interfaces/ --go-grpc_out=./interfaces/ -I ./contracts/proto/ ./contracts/proto/*/**/*.proto

# Generate mock implementations
generate-mocks:
	@echo Generating Mock RPC Clients
	@go run github.com/golang/mock/mockgen github.com/nitrictech/go-sdk/interfaces/nitric/v1 DocumentServiceClient,EventServiceClient,TopicServiceClient,QueueServiceClient,StorageServiceClient,FaasServiceClient,FaasService_TriggerStreamClient,DocumentService_QueryStreamClient,SecretServiceClient > mocks/clients.go

# Runs tests for coverage upload to codecov.io
test-ci: generate-mocks
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./...

test: generate-mocks
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover ./...