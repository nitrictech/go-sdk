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

package apis

type (
	ApiOption    func(api *api)
	RouteOption  func(route Route)
	MethodOption func(mo *methodOptions)
)

type JwtSecurityRule struct {
	Issuer    string
	Audiences []string
}

type methodOptions struct {
	security         []OidcOptions
	securityDisabled bool
}

// WithMiddleware - Apply a middleware function to all handlers in the API
func WithMiddleware(middleware Middleware) ApiOption {
	return func(api *api) {
		api.middleware = middleware
	}
}

func OidcRule(name string, issuer string, audiences []string) SecurityOption {
	return func(scopes []string) OidcOptions {
		return OidcOptions{
			Name:      name,
			Issuer:    issuer,
			Audiences: audiences,
			Scopes:    scopes,
		}
	}
}

// WithSecurityJwtRule - Apply a JWT security rule to the API
func WithSecurityJwtRule(name string, rule JwtSecurityRule) ApiOption {
	return func(api *api) {
		if api.securityRules == nil {
			api.securityRules = make(map[string]interface{})
		}

		api.securityRules[name] = rule
	}
}

// WithSecurity - Apply security settings to the API
func WithSecurity(oidcOptions OidcOptions) ApiOption {
	return func(api *api) {
		if api.security == nil {
			api.security = []OidcOptions{oidcOptions}
		} else {
			api.security = append(api.security, oidcOptions)
		}
	}
}

// WithPath - Set the base path for the API
func WithPath(path string) ApiOption {
	return func(api *api) {
		api.path = path
	}
}

// WithNoMethodSecurity - Disable security for a method
func WithNoMethodSecurity() MethodOption {
	return func(mo *methodOptions) {
		mo.securityDisabled = true
	}
}

// WithMethodSecurity - Override/set the security settings for a method
func WithMethodSecurity(oidcOptions OidcOptions) MethodOption {
	return func(mo *methodOptions) {
		mo.securityDisabled = false
		if mo.security == nil {
			mo.security = []OidcOptions{oidcOptions}
		} else {
			mo.security = append(mo.security, oidcOptions)
		}
	}
}
