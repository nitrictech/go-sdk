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

package storage

import (
	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

// Storage - Nitric storage API client
type Storage interface {
	// Bucket - Get a bucket reference for the provided name
	Bucket(name string) Bucket
}

type storageImpl struct {
	storageClient v1.StorageClient
}

func (s *storageImpl) Bucket(name string) Bucket {
	return &bucketImpl{
		storageClient: s.storageClient,
		name:          name,
	}
}

// New - Create a new Storage client with default options
func New() (Storage, error) {
	conn, err := grpc.NewClient(constants.NitricAddress(), constants.DefaultOptions()...)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"Storage.New: Unable to reach StorageServiceServer",
			err,
		)
	}

	sClient := v1.NewStorageClient(conn)

	return &storageImpl{
		storageClient: sClient,
	}, nil
}
