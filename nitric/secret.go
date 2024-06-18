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

package nitric

import (
	"context"
	"fmt"

	"github.com/nitrictech/go-sdk/api/secrets"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type SecretPermission string

const (
	SecretAccess SecretPermission = "access"
	SecretPut    SecretPermission = "put"
)

var SecretEverything []SecretPermission = []SecretPermission{SecretAccess, SecretPut}

type Secret interface{}

type secret struct {
	name    string
	manager Manager
}

func NewSecret(name string) *secret {
	return &secret{
		name:    name,
		manager: defaultManager,
	}
}

func (s *secret) Allow(permission SecretPermission, permissions ...SecretPermission) (secrets.SecretRef, error) {
	allPerms := append([]SecretPermission{permission}, permissions...)

	return defaultManager.newSecret(s.name, allPerms...)
}

func (m *manager) newSecret(name string, permissions ...SecretPermission) (secrets.SecretRef, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	colRes := &v1.ResourceIdentifier{
		Type: v1.ResourceType_Secret,
		Name: name,
	}

	dr := &v1.ResourceDeclareRequest{
		Id: colRes,
		Config: &v1.ResourceDeclareRequest_Secret{
			Secret: &v1.SecretResource{},
		},
	}
	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []v1.Action{}
	for _, perm := range permissions {
		switch perm {
		case SecretAccess:
			actions = append(actions, v1.Action_SecretAccess)
		case SecretPut:
			actions = append(actions, v1.Action_SecretPut)
		default:
			return nil, fmt.Errorf("secretPermission %s unknown", perm)
		}
	}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(colRes, actions))
	if err != nil {
		return nil, err
	}

	if m.secrets == nil {
		m.secrets, err = secrets.New()
		if err != nil {
			return nil, err
		}
	}

	return m.secrets.Secret(name), nil
}
