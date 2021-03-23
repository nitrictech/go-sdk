install:
	@echo installing go dependencies
	@go mod download

install-tools: install
	@echo Installing tools from tools.go
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

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