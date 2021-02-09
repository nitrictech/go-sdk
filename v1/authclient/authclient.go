package authclient

import (
	"context"
	"fmt"
	v1 "go.nitric.io/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

type AuthClient interface {
	CreateUser(tenant string, userId string, email string, password string) error
}

// NitricAuthClient - gRPC based client to nitric membrane server for auth services.
type NitricAuthClient struct {
	conn *grpc.ClientConn
	c v1.AuthClient
}

// CreateUser - create a new user in the provided specific auth service.
func (a NitricAuthClient) CreateUser(tenant string, userId string, email string, password string) error {
	_, err := a.c.CreateUser(context.Background(), &v1.CreateUserRequest{
		Tenant:   tenant,
		Id:       userId,
		Email:    email,
		Password: password,
	})

	// TODO: Should something be returned?
	return err
}

// Close - closes the connection to the membrane server
// no need to call close if the connect is to remain open for the lifetime of the application.
func (a NitricAuthClient) Close() error {
	return a.conn.Close()
}

// FIXME: Extract into shared code.
// NewAuthClient - create a new nitric auth client
func NewAuthClient() (AuthClient, error) {
	// Connect to the gRPC Membrane Server
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to Membrane gRPC server: %s", err)
	}

	return &NitricAuthClient{
		conn: conn,
		c: v1.NewAuthClient(conn),
	}, nil
}