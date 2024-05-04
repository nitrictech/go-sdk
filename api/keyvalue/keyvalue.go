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

package keyvalue

import (
	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
)

// KeyValue - Idiomatic interface for the nitric Key Value Store Service
type KeyValue interface {
	// Gets a store instance that refers to the store at the specified path.
	Store(string) Store
}

type keyValueImpl struct {
	kvClient v1.KvStoreClient
}

func (k *keyValueImpl) Store(name string) Store {
	return &storeImpl{
		name:     name,
		kvClient: k.kvClient,
	}
}

// New - Construct a new Key Value Store Client with default options
func New() (KeyValue, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"KeyValue.New: Unable to reach KVStoreServiceServer",
			err,
		)
	}

	kvClient := v1.NewKvStoreClient(conn)

	return &keyValueImpl{
		kvClient: kvClient,
	}, nil
}
