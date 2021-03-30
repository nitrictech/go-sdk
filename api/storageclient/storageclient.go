package storageclient

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

type StorageClient interface {
	Read(bucketName string, key string) ([]byte, error)
	Write(bucketName string, key string, body []byte) error
}

type NitricStorageClient struct {
	conn *grpc.ClientConn
	c    v1.StorageClient
}

// Get - retrieves an exist item from a bucket by its key
func (s NitricStorageClient) Read(bucketName string, key string) ([]byte, error) {
	res, err := s.c.Read(context.Background(), &v1.StorageReadRequest{
		BucketName: bucketName,
		Key:        key,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get content with key [%s] from bucket [%s]: %s", key, bucketName, err)
	}

	return res.GetBody(), nil
}

// Put - stores an item in a bucket under the given key.
func (s NitricStorageClient) Write(bucketName string, key string, body []byte) error {
	_, err := s.c.Write(context.Background(), &v1.StorageWriteRequest{
		BucketName: bucketName,
		Key:        key,
		Body:       body,
	})

	return err
}

func NewStorageClient(conn *grpc.ClientConn) StorageClient {
	return &NitricStorageClient{
		conn: conn,
		c:    v1.NewStorageClient(conn),
	}
}

func NewWithClient(client v1.StorageClient) StorageClient {
	return &NitricStorageClient{
		c: client,
	}
}
