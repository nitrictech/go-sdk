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

package resources

type ApiOption = func(api *api)

type JwtSecurityRule struct {
	Issuer    string
	Audiences []string
}

func WithSecurityJwtRule(name string, rule JwtSecurityRule) ApiOption {
	return func(api *api) {
		if api.securityRules == nil {
			api.securityRules = make(map[string]interface{})
		}

		api.securityRules[name] = rule
	}
}

func WithSecurity(name string, scopes []string) ApiOption {
	return func(api *api) {
		if api.security == nil {
			api.security = make(map[string][]string)
		}

		api.security[name] = scopes
	}
}
