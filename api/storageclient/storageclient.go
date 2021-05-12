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

package storageclient

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

type ReadOptions struct {
	Bucket string
	Key    string
}

type ReadResult struct {
	Data []byte
}

type WriteOptions struct {
	Bucket string
	Key    string
	Data   []byte
}

// Empty response for forwards compatibility
type WriteResult struct{}

type DeleteOptions struct {
	Bucket string
	Key    string
}

// Empty response for forwards compatibility
type DeleteResult struct{}

type StorageClient interface {
	Read(*ReadOptions) (*ReadResult, error)
	Write(*WriteOptions) (*WriteResult, error)
	Delete(*DeleteOptions) (*DeleteResult, error)
}

type NitricStorageClient struct {
	conn *grpc.ClientConn
	c    v1.StorageClient
}

// Get - retrieves an exist item from a bucket by its key
func (s *NitricStorageClient) Read(opts *ReadOptions) (*ReadResult, error) {
	res, err := s.c.Read(context.Background(), &v1.StorageReadRequest{
		BucketName: opts.Bucket,
		Key:        opts.Key,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get content with key [%s] from bucket [%s]: %s", opts.Key, opts.Bucket, err)
	}

	return &ReadResult{
		Data: res.GetBody(),
	}, nil
}

// Put - stores an item in a bucket under the given key.
func (s *NitricStorageClient) Write(opts *WriteOptions) (*WriteResult, error) {
	if _, err := s.c.Write(context.Background(), &v1.StorageWriteRequest{
		BucketName: opts.Bucket,
		Key:        opts.Key,
		Body:       opts.Data,
	}); err != nil {
		return nil, err
	}

	return &WriteResult{}, nil
}

func (s *NitricStorageClient) Delete(opts *DeleteOptions) (*DeleteResult, error) {
	if _, err := s.c.Delete(context.Background(), &v1.StorageDeleteRequest{
		BucketName: opts.Bucket,
		Key:        opts.Key,
	}); err != nil {
		return nil, err
	}

	return &DeleteResult{}, nil
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
