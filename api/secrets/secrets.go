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
	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

// Secrets - Base client for the Nitric Secrets service
type Secrets interface {
	// Secret - Creates a new secret reference
	Secret(string) SecretRef
}

type secretsImpl struct {
	secretClient v1.SecretServiceClient
}

func (s *secretsImpl) Secret(name string) SecretRef {
	return &secretRefImpl{
		name:         name,
		secretClient: s.secretClient,
	}
}

// New - Create a new Secrets client
func New() (Secrets, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"Secrets.New: Unable to reach SecretsServiceServer",
			err,
		)
	}

	sClient := v1.NewSecretServiceClient(conn)

	return &secretsImpl{
		secretClient: sClient,
	}, nil
}
