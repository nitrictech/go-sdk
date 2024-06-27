// Copyright 2023 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keyvalue

import (
	"context"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/protoutils"

	v1 "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
)

type ScanKeysRequest = v1.KvStoreScanKeysRequest

type ScanKeysOption = func(*ScanKeysRequest)

// Apply a prefix to the scan keys request
func WithPrefix(prefix string) ScanKeysOption {
	return func(req *ScanKeysRequest) {
		req.Prefix = prefix
	}
}

// TODO: maybe move keystream to separate file
type KeyStream struct {
	stream v1.KvStore_ScanKeysClient
}

func (k *KeyStream) Recv() (string, error) {
	resp, err := k.stream.Recv()
	if err != nil {
		return "", err
	}

	return resp.Key, nil
}

type Store interface {
	// Name - The name of the store
	Name() string
	// Get a value from the store
	Get(context.Context, string) (map[string]interface{}, error)
	// Set a value in the store
	Set(context.Context, string, map[string]interface{}) error
	// Delete a value from the store
	Delete(context.Context, string) error
	// Return an async iterable of keys in the store
	Keys(context.Context, ...ScanKeysOption) (*KeyStream, error)
}

type storeImpl struct {
	name     string
	kvClient v1.KvStoreClient
}

func (s *storeImpl) Name() string {
	return s.name
}

func (s *storeImpl) Get(ctx context.Context, key string) (map[string]interface{}, error) {
	ref := &v1.ValueRef{
		Store: s.name,
		Key:   key,
	}

	r, err := s.kvClient.GetValue(ctx, &v1.KvStoreGetValueRequest{
		Ref: ref,
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	val := r.GetValue()
	if val == nil {
		return nil, errors.New(codes.NotFound, "Key not found")
	}
	content := val.GetContent().AsMap()

	return content, nil
}

func (s *storeImpl) Set(ctx context.Context, key string, value map[string]interface{}) error {
	ref := &v1.ValueRef{
		Store: s.name,
		Key:   key,
	}

	// Convert payload to Protobuf Struct
	contentStruct, err := protoutils.NewStruct(value)
	if err != nil {
		return errors.NewWithCause(codes.InvalidArgument, "Store.Set", err)
	}

	_, err = s.kvClient.SetValue(ctx, &v1.KvStoreSetValueRequest{
		Ref:     ref,
		Content: contentStruct,
	})
	if err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

func (s *storeImpl) Delete(ctx context.Context, key string) error {
	ref := &v1.ValueRef{
		Store: s.name,
		Key:   key,
	}

	_, err := s.kvClient.DeleteKey(ctx, &v1.KvStoreDeleteKeyRequest{
		Ref: ref,
	})
	if err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

func (s *storeImpl) Keys(ctx context.Context, opts ...ScanKeysOption) (*KeyStream, error) {
	store := &v1.Store{
		Name: s.name,
	}

	request := &v1.KvStoreScanKeysRequest{
		Store:  store,
		Prefix: "",
	}

	// Apply options to the request payload - Prefix modification
	for _, opt := range opts {
		opt(request)
	}

	streamClient, err := s.kvClient.ScanKeys(ctx, request)
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	return &KeyStream{
		stream: streamClient,
	}, nil
}
