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
	"fmt"

	"github.com/nitrictech/go-sdk/nitric/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type SecretPermission string

const (
	SecretAccess SecretPermission = "access"
	SecretPut    SecretPermission = "put"
)

var SecretEverything []SecretPermission = []SecretPermission{SecretAccess, SecretPut}

type Secret interface {
	// Allow requests the given permissions to the secret.
	Allow(permission SecretPermission, permissions ...SecretPermission) *SecretClient
}

type secret struct {
	name         string
	manager      *workers.Manager
	registerChan <-chan workers.RegisterResult
}

// NewSecret - Create a new Secret resource
func NewSecret(name string) *secret {
	secret := &secret{
		name:    name,
		manager: workers.GetDefaultManager(),
	}

	secret.registerChan = secret.manager.RegisterResource(&v1.ResourceDeclareRequest{
		Id: &v1.ResourceIdentifier{
			Type: v1.ResourceType_Secret,
			Name: name,
		},
		Config: &v1.ResourceDeclareRequest_Secret{
			Secret: &v1.SecretResource{},
		},
	})

	return secret
}

func (s *secret) Allow(permission SecretPermission, permissions ...SecretPermission) *SecretClient {
	allPerms := append([]SecretPermission{permission}, permissions...)

	actions := []v1.Action{}
	for _, perm := range allPerms {
		switch perm {
		case SecretAccess:
			actions = append(actions, v1.Action_SecretAccess)
		case SecretPut:
			actions = append(actions, v1.Action_SecretPut)
		default:
			panic(fmt.Sprintf("secretPermission %s unknown", perm))
		}
	}

	registerResult := <-s.registerChan
	if registerResult.Err != nil {
		panic(registerResult.Err)
	}

	err := s.manager.RegisterPolicy(registerResult.Identifier, actions...)
	if err != nil {
		panic(err)
	}

	client, err := NewSecretClient(s.name)
	if err != nil {
		panic(err)
	}

	return client
}
