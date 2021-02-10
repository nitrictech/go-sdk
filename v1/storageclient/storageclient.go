package storageclient

import (
	"context"
	"fmt"

	v1 "go.nitric.io/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

type StorageClient interface {
	Get(bucketName string, key string) ([]byte, error)
	Put(bucketName string, key string, body []byte) error
}

type NitricStorageClient struct {
	conn *grpc.ClientConn
	c    v1.StorageClient
}

// Get - retrieves an exist item from a bucket by its key
func (s NitricStorageClient) Get(bucketName string, key string) ([]byte, error) {
	res, err := s.c.Get(context.Background(), &v1.GetRequest{
		BucketName: bucketName,
		Key:        key,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get content with key [%s] from bucket [%s]: %s", key, bucketName, err)
	}

	return res.GetBody(), nil
}

// Put - stores an item in a bucket under the given key.
func (s NitricStorageClient) Put(bucketName string, key string, body []byte) error {
	res, err := s.c.Put(context.Background(), &v1.PutRequest{
		BucketName: bucketName,
		Key:        key,
		Body:       body,
	})

	// FIXME: we probably shouldn't return a success boolean. Errors should indicate a failure.
	if res != nil && !res.GetSuccess() {
		return fmt.Errorf("failed to store data, unexpected failure status returned")
	}

	return err
}

func NewStorageClient(conn *grpc.ClientConn) StorageClient {
	return &NitricStorageClient{
		conn: conn,
		c:    v1.NewStorageClient(conn),
	}
}
