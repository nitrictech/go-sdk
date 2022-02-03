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

	v1 "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
)

// Storage - Nitric storage API client
type Storage interface {
	// Bucket - Get a bucket reference for the provided name
	Bucket(name string) Bucket
}

type storageImpl struct {
	sc v1.StorageServiceClient
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
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"Storage.New: Unable to reach StorageServiceServer",
			err,
		)
	}

	sClient := v1.NewStorageServiceClient(conn)

	return &storageImpl{
		sc: sClient,
	}, nil
}
