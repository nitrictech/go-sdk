// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"fmt"

	"github.com/nitrictech/go-sdk/api/eventclient"
	"github.com/nitrictech/go-sdk/api/kvclient"
	"github.com/nitrictech/go-sdk/api/queueclient"
	"github.com/nitrictech/go-sdk/api/storageclient"
	"github.com/nitrictech/go-sdk/constants"
	"google.golang.org/grpc"
)

// NitricClient - provider services client
// TODO: Look at adding generics for scope: https://blog.golang.org/generics-next-step
type NitricClient interface {
	KV() kvclient.KVClient
	Event() eventclient.EventClient
	Queue() queueclient.QueueClient
	Storage() storageclient.StorageClient
	Close()
}

// Client - NitricClient gRPC implementation
type Client struct {
	connection *grpc.ClientConn
	kv         kvclient.KVClient
	eventing   eventclient.EventClient
	queue      queueclient.QueueClient
	storage    storageclient.StorageClient
}

// Re-connect a previously closed client
func (c *Client) ensureGrpcConnection() {
	if c.connection == nil {
		// Create a new connection
		// Let clients naturally return their errors?
		conn, _ := grpc.Dial(constants.NitricAddress(), grpc.WithInsecure())
		c.connection = conn
	}
}

// KV - returns a kv service client
func (c *Client) KV() kvclient.KVClient {
	c.ensureGrpcConnection()
	if c.kv == nil {
		c.kv = kvclient.NewKVClient(c.connection)
	}

	return c.kv
}

// Storage - returns a storage client
func (c *Client) Storage() storageclient.StorageClient {
	c.ensureGrpcConnection()
	if c.storage == nil {
		c.storage = storageclient.NewStorageClient(c.connection)
	}

	return c.storage
}

// Event - returns an event service client
func (c *Client) Event() eventclient.EventClient {
	c.ensureGrpcConnection()
	if c.eventing == nil {
		c.eventing = eventclient.NewEventClient(c.connection)
	}

	return c.eventing
}

// Queue - returns a queue service client
func (c *Client) Queue() queueclient.QueueClient {
	c.ensureGrpcConnection()
	if c.queue == nil {
		c.queue = queueclient.NewQueueClient(c.connection)
	}

	return c.queue
}

// Close - close the nitric client, its connection to the nitric service and all service clients
func (c *Client) Close() {
	if c.connection != nil {
		_ = c.connection.Close()
		// Nil out existing clients
		c.queue = nil
		c.eventing = nil
		c.storage = nil
		c.kv = nil
	}
}

// New - constructs a new NitricClient with a connection to the nitric service.
// connection information is retrieved from the standard environment variable
func New() (NitricClient, error) {
	// Connect to the gRPC Membrane Server
	conn, err := grpc.Dial(contants.NitricAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to Membrane gRPC server: %s", err)
	}

	return &Client{
		connection: conn,
	}, nil
}
