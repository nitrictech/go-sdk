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

// SecretVersionRef - A reference to a secret version
type SecretVersionRef interface {
	// Access - Retrieve the value of the secret
	Access(ctx context.Context) (SecretValue, error)
	Secret() SecretRef
	Version() string
}

type secretVersionRefImpl struct {
	secretClient v1.SecretServiceClient
	secret       SecretRef
	version      string
}

func (s *secretVersionRefImpl) Secret() SecretRef {
	return s.secret
}

func (s *secretVersionRefImpl) Version() string {
	return s.version
}

func (s *secretVersionRefImpl) Access(ctx context.Context) (SecretValue, error) {
	r, err := s.secretClient.Access(ctx, &v1.SecretAccessRequest{
		SecretVersion: &v1.SecretVersion{
			Secret: &v1.Secret{
				Name: s.secret.Name(),
			},
			Version: s.version,
		},
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	return &secretValueImpl{
		version: &secretVersionRefImpl{
			secretClient: s.secretClient,
			secret:       s.secret,
			version:      r.GetSecretVersion().GetVersion(),
		},
		val: r.GetValue(),
	}, nil
}
