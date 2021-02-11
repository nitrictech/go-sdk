package authclient

import (
	"context"
	v1 "go.nitric.io/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

type AuthClient interface {
	CreateUser(tenant string, userId string, email string, password string) error
}

// NitricAuthClient - gRPC based client to nitric membrane server for auth services.
type NitricAuthClient struct {
	conn *grpc.ClientConn
	c    v1.AuthClient
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

// FIXME: Extract into shared code.
// NewAuthClient - create a new nitric auth client
func NewAuthClient(conn *grpc.ClientConn) AuthClient {
	return &NitricAuthClient{
		conn: conn,
		c:    v1.NewAuthClient(conn),
	}
}

func NewWithClient(client v1.AuthClient) AuthClient {
	return &NitricAuthClient{
		c:    client,
	}
}
