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
	"strings"

	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type OidcOptions struct {
	Name    string
	Issuer 	string
	Audiences []string
	Scopes []string
}

func attachOidc(apiName string, options OidcOptions) error {
	_, err := NewOidcSecurityDefinition(apiName, options)
	if err != nil {
		return err
	}
	return nil
}

type SecurityOption = func (scopes []string) OidcOptions

type OidcSecurityDefinition interface {

}

type oidcSecurityDefinition struct {
	OidcSecurityDefinition

	ApiName    string
	RuleName	string
	Issuer		string
	Audiences	[]string

	manager Manager
}

func NewOidcSecurityDefinition(apiName string, options OidcOptions) (OidcSecurityDefinition, error) {
	return defaultManager.newOidcSecurityDefinition(apiName, options)
}

func (m *manager) newOidcSecurityDefinition(apiName string, options OidcOptions) (OidcSecurityDefinition, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	o := &oidcSecurityDefinition{
		ApiName: apiName,
		RuleName: options.Name,
		Issuer: options.Issuer,
		Audiences: options.Audiences,
		manager: m,
	}

	// declare resource
	_, err = rsc.Declare(context.TODO(), &v1.ResourceDeclareRequest{
		Id: &v1.ResourceIdentifier{
			Name: strings.Join([]string{
				options.Name,
				apiName,
			}, "-"),
			Type: v1.ResourceType_ApiSecurityDefinition,
		},
		Config: &v1.ResourceDeclareRequest_ApiSecurityDefinition{
			ApiSecurityDefinition: &v1.ApiSecurityDefinitionResource{
				ApiName: apiName,
				Definition: &v1.ApiSecurityDefinitionResource_Oidc{
					Oidc: &v1.ApiOpenIdConnectionDefinition{
						Issuer: o.Issuer,
						Audiences: o.Audiences,
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return o, nil
}

