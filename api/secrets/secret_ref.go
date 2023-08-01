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

package secrets

import (
	"context"

	"github.com/nitrictech/go-sdk/api/errors"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

// SecretRef is a reference to a cloud secret for secret storage.
type SecretRef interface {
	Name() string
	Put(context.Context, []byte) (SecretVersionRef, error)
	Version(string) SecretVersionRef
	Latest() SecretVersionRef
}

type secretRefImpl struct {
	name string
	secretClient   v1.SecretServiceClient
}

func (s *secretRefImpl) Name() string {
	return s.name
}

func (s *secretRefImpl) Put(ctx context.Context, sec []byte) (SecretVersionRef, error) {
	resp, err := s.secretClient.Put(ctx, &v1.SecretPutRequest{
		Secret: &v1.Secret{
			Name: s.name,
		},
		Value: sec,
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	return &secretVersionRefImpl{
		secretClient:      s.secretClient,
		version: resp.GetSecretVersion().Version,
		secret:  s,
	}, nil
}

func (s *secretRefImpl) Version(name string) SecretVersionRef {
	return &secretVersionRefImpl{
		secret:  s,
		secretClient:      s.secretClient,
		version: name,
	}
}

func (s *secretRefImpl) Latest() SecretVersionRef {
	return s.Version("latest")
}
