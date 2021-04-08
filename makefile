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

license-check: install-tools
	@echo Checking OSS Licenses
	@go build -o ./bin/licenses ./licenses.go 
	@lichen --config=./lichen.yaml ./bin/licenses

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
	@go run github.com/golang/mock/mockgen github.com/nitrictech/go-sdk/interfaces/nitric/v1 UserClient,KeyValueClient,EventClient,TopicClient,QueueClient,StorageClient  > mocks/clients.go

test: generate-mocks
	@echo Testing Nitric Go SDK
	@go run github.com/onsi/ginkgo/ginkgo -cover ./api/...