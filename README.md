<p align="center">
  <img src="./docs/assets/dot_matrix_logo_go.png" alt="Nitric Logo"/>
</p>

# Nitric SDK for Go

![test status](https://github.com/nitrictech/go-sdk/actions/workflows/test.yaml/badge.svg?branch=main)

Client libarary for interfacing with the Nitric APIs as well as the creation of golang functions with Nitric.

## Quick Start

### Using the Nitric CLI

1. Get the Nitric CLI
2. Create a new Nitric Project `nitric make:project <my-new-project>`
3. Select `function/golang15` as a starter function

## Using the Nitric SDK

### Creating a new API client
```go
import "github.com/nitrictech/go-sdk/api"

// NitricFunction - Handles individual function requests (http, events, etc.)
func createNitricClient() {
  client, err := api.New()

  if err != nil {
    // Do something with err
  }
}
```

### Starting a Nitric FaaS server

```go
package main

import "github.com/nitrictech/go-sdk/faas"

// NitricFunction - Handles individual function requests (http, events, etc.)
func NitricFunction(request *faas.NitricRequest) *faas.NitricResponse {
	// Do something interesting...
	return &faas.NitricResponse{
		Status: 200,
		Body:   []byte("Hello Nitric"),
	}
}

func main() {
	faas.Start(NitricFunction)
}
```
<!-- TODO: Add additional examples but don't add too much noise to the landing README -->
<!-- More specific usage examples can be found in [examples](./examples/README.md) -->

