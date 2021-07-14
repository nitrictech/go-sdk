package storage

import (
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

// Storage - Nitric storage API client
type Storage interface {
	// Bucket - Get a bucket reference for the provided name
	Bucket(name string) Bucket
}

type storageImpl struct {
	sc v1.StorageClient
}

func (s *storageImpl) Bucket(name string) Bucket {
	return &bucketImpl{
		sc:   s.sc,
		name: name,
	}
}

// New - Create a new Storage client with default options
func New() (Storage, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)

	if err != nil {
		return nil, err
	}

	sClient := v1.NewStorageClient(conn)

	return &storageImpl{
		sc: sClient,
	}, nil
}
