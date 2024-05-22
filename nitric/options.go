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
	"github.com/nitrictech/go-sdk/handler"
)

type (
	ApiOption    = func(api *api)
	MethodOption = func(mo *methodOptions)
)

type JwtSecurityRule struct {
	Issuer    string
	Audiences []string
}

type methodOptions struct {
	security         map[string][]string
	securityDisabled bool
}

func OidcRule(name string,  issuer string , audiences []string) SecurityOption {
	return func(scopes []string) OidcOptions {
		return OidcOptions{
			Name: name,
			Issuer: issuer,
			Audiences: audiences,
			Scopes: scopes,
		}
	}
}

func WithMiddleware(middleware ...handler.HttpMiddleware) ApiOption {
	return func(api *api) {
		if api.middleware != nil {
			api.middleware = handler.ComposeHttpMiddleware(api.middleware, handler.ComposeHttpMiddleware(middleware...))
		} else {
			api.middleware = handler.ComposeHttpMiddleware(middleware...)
		}
	}
}

func WithSecurityJwtRule(name string, rule JwtSecurityRule) ApiOption {
	return func(api *api) {
		if api.securityRules == nil {
			api.securityRules = make(map[string]interface{})
		}

		api.securityRules[name] = rule
	}
}

func WithSecurity(oidcOptions OidcOptions) ApiOption {
	return func(api *api) {
		if api.security == nil {
			api.security = []OidcOptions{oidcOptions}
		} else {
			api.security = append(api.security, oidcOptions)
		}
	}	
}

// WithPath - Prefixes API with the given path
func WithPath(path string) ApiOption {
	return func(api *api) {
		api.path = path
	}
}

func WithNoMethodSecurity() MethodOption {
	return func(mo *methodOptions) {
		mo.securityDisabled = true
	}
}

func WithMethodSecurity(name string, scopes []string) MethodOption {
	return func(mo *methodOptions) {
		if name == "" {
			mo.securityDisabled = true
		} else {
			if mo.security == nil {
				mo.security = map[string][]string{}
			}

			mo.security[name] = scopes
		}
	}
}
