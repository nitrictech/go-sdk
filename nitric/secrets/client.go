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

	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
)

type SecretValue []byte

func (s SecretValue) AsString() string {
	return string(s)
}

type SecretClientIface interface {
	// Name - Return the name of this secret
	Name() string
	// Put - Store a new value in this secret, returning a reference to the new version created
	Put(context.Context, []byte) (string, error)
	// Access - Access the latest version of this secret
	Access(context.Context) (SecretValue, error)
	// AccessVersion - Access a specific version of the secret
	AccessVersion(context.Context, string) (SecretValue, error)
}

var _ SecretClientIface = (*SecretClient)(nil)

// SecretClient - Reference to a cloud secret
type SecretClient struct {
	name         string
	secretClient v1.SecretManagerClient
}

// Name - Return the name of this secret
func (s *SecretClient) Name() string {
	return s.name
}

// Put - Store a new value in this secret, returning a reference to the new version created
func (s *SecretClient) Put(ctx context.Context, sec []byte) (string, error) {
	resp, err := s.secretClient.Put(ctx, &v1.SecretPutRequest{
		Secret: &v1.Secret{
			Name: s.name,
		},
		Value: sec,
	})
	if err != nil {
		return "", errors.FromGrpcError(err)
	}

	return resp.GetSecretVersion().Version, nil
}

const latestVersionId = "latest"

// Access - Access the latest version of this secret
func (s *SecretClient) Access(ctx context.Context) (SecretValue, error) {
	return s.AccessVersion(ctx, latestVersionId)
}

// AccessVersion - Access a specific version of the secret
func (s *SecretClient) AccessVersion(ctx context.Context, version string) (SecretValue, error) {
	r, err := s.secretClient.Access(ctx, &v1.SecretAccessRequest{
		SecretVersion: &v1.SecretVersion{
			Secret: &v1.Secret{
				Name: s.name,
			},
			Version: version,
		},
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	return SecretValue(r.GetValue()), nil
}

func NewSecretClient(name string) (*SecretClient, error) {
	conn, err := grpc.NewClient(constants.NitricAddress(), constants.DefaultOptions()...)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"NewSecretClient: unable to reach nitric server",
			err,
		)
	}

	sClient := v1.NewSecretManagerClient(conn)

	return &SecretClient{
		secretClient: sClient,
		name:         name,
	}, nil
}
