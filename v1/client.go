package v1

import (
	"fmt"
	"net"
	"os"

	"go.nitric.io/go-sdk/v1/authclient"
	"go.nitric.io/go-sdk/v1/documentsclient"
	"go.nitric.io/go-sdk/v1/eventclient"
	"go.nitric.io/go-sdk/v1/queueclient"
	"go.nitric.io/go-sdk/v1/storageclient"
	"google.golang.org/grpc"
)

const (
	nitricServiceHostDefault    = "127.0.0.1"
	nitricServicePortDefault    = "50051"
	nitricServiceHostEnvVarName = "NITRIC_SERVICE_HOST"
	nitricServicePortEnvVarName = "NITRIC_SERVICE_PORT"
)

// NitricClient - provider services client
// TODO: Look at adding generics for scope: https://blog.golang.org/generics-next-step
type NitricClient interface {
	Auth() authclient.AuthClient
	Documents() documentsclient.DocumentsClient
	Eventing() eventclient.EventClient
	Queue() queueclient.QueueClient
	Storage() storageclient.StorageClient
	Close()
}

// Client - NitricClient gRPC implementation
type Client struct {
	connection *grpc.ClientConn
	auth       authclient.AuthClient
	documents  documentsclient.DocumentsClient
	eventing   eventclient.EventClient
	queue      queueclient.QueueClient
	storage    storageclient.StorageClient
}

// Re-connect a previously closed client
func (c *Client) ensureGrpcConnection() {
	if c.connection == nil {
		// Create a new connection
		// Let clients naturally return their errors?
		conn, _ := grpc.Dial(nitricAddress(), grpc.WithInsecure())
		c.connection = conn
	}
}

// Auth - returns an auth service client
func (c *Client) Auth() authclient.AuthClient {
	c.ensureGrpcConnection()
	if c.auth == nil {
		c.auth = authclient.NewAuthClient(c.connection)
	}

	return c.auth
}

// Documents - returns a document service client
func (c *Client) Documents() documentsclient.DocumentsClient {
	c.ensureGrpcConnection()
	if c.documents == nil {
		c.documents = documentsclient.NewDocumentsClient(c.connection)
	}

	return c.documents
}

// Storage - returns a storage client
func (c *Client) Storage() storageclient.StorageClient {
	c.ensureGrpcConnection()
	if c.storage == nil {
		c.storage = storageclient.NewStorageClient(c.connection)
	}

	return c.storage
}

// Eventing - retuns an eventing service client
func (c *Client) Eventing() eventclient.EventClient {
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
		c.auth = nil
		c.documents = nil
	}
}

// getEnvWithFallback - Returns an envirable variable's value from its name or a default value if the variable isn't set
func getEnvWithFallback(varName string, fallback string) string {
	if v := os.Getenv(varName); v == "" {
		return fallback
	} else {
		return v
	}
}

// nitricPort - retrieves the environment variable which specifies the port (e.g. 50051) of the nitric service
//
// if the env var isn't set, returns the default port
func nitricPort() string {
	return getEnvWithFallback(nitricServicePortEnvVarName, nitricServicePortDefault)
}

// nitricPort - retrieves the environment variable which specifies the host (e.g. 127.0.0.1) of the nitric service
//
// if the env var isn't set, returns the default host
func nitricHost() string {
	return getEnvWithFallback(nitricServiceHostEnvVarName, nitricServiceHostDefault)
}

// nitricAddress - constructs the full address i.e. host:port, of the nitric service based on config or defaults
func nitricAddress() string {
	return net.JoinHostPort(nitricAddress(), nitricPort())
}

// New - constructs a new NitricClient with a connection to the nitric service.
// connection information is retrieved from the standard environment variable
func New() (NitricClient, error) {
	// Connect to the gRPC Membrane Server
	conn, err := grpc.Dial(nitricAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to Membrane gRPC server: %s", err)
	}

	return &Client{
		connection: conn,
	}, nil
}
