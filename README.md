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

### Creating a new API clients
```go
import "github.com/nitrictech/go-sdk/api/events"

// NitricFunction - Handles individual function requests (http, events, etc.)
func createNitricClient() {
	ec, err := events.New()

  if err != nil {
    // Do something with err
  }

	// use the new events client
	// ec.Topic("my-topic").Publish(...)
}
```

### Starting a Nitric FaaS server

```go
package main

import "github.com/nitrictech/go-sdk/faas"

// NitricFunction - Handles individual function requests (http, events, etc.)
func NitricFunction(trigger *faas.NitricTrigger) (*faas.TriggerResponse, error) {
	// Construct a default response base on the existing trigger context
	resp := trigger.DefaultResponse()
	resp.SetData([]byte("Hello Nitric"))

	return resp, nil
}

func main() {
	faas.Start(NitricFunction)
}
```
<!-- TODO: Add additional examples but don't add too much noise to the landing README -->
<!-- More specific usage examples can be found in [examples](./examples/README.md) -->

