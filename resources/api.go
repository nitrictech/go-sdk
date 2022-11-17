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

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/nitrictech/go-sdk/faas"
	v1 "github.com/nitrictech/go-sdk/nitric/v1"
)

// Route providers convenience functions to register a handler in a single method.
type Route interface {
	Get(handler faas.HttpMiddleware, opts ...MethodOption)
	Patch(handler faas.HttpMiddleware, opts ...MethodOption)
	Put(handler faas.HttpMiddleware, opts ...MethodOption)
	Post(handler faas.HttpMiddleware, opts ...MethodOption)
	Delete(handler faas.HttpMiddleware, opts ...MethodOption)
	Options(handler faas.HttpMiddleware, opts ...MethodOption)
}

type route struct {
	apiPath string
	apiName string
	m       *manager
}

func NewRoute(apiName, apiPath string) Route {
	return run.NewRoute(apiName, apiPath)
}

func (m *manager) NewRoute(apiName, apiPath string) Route {
	return &route{
		m:       m,
		apiPath: apiPath,
		apiName: apiName,
	}
}

func (r *route) AddMethodHandler(methods []string, handler faas.HttpMiddleware, opts ...MethodOption) {
	bName := path.Join(r.apiName, r.apiPath, strings.Join(methods, "-"))

	mo := &methodOptions{}
	for _, o := range opts {
		o(mo)
	}

	_, ok := r.m.builders[bName]
	if !ok {
		r.m.builders[bName] = faas.New().WithApiWorkerOpts(faas.ApiWorkerOptions{
			ApiName:          r.apiName,
			Path:             r.apiPath,
			Security:         mo.security,
			SecurityDisabled: mo.securityDisabled,
		})
	}

	b := r.m.builders[bName]

	for _, m := range methods {
		b.Http(m, handler)
	}

	r.m.addStarter("route:"+bName, b)
	r.m.builders[bName] = b
}

func (r *route) Get(handler faas.HttpMiddleware, opts ...MethodOption) {
	r.AddMethodHandler([]string{http.MethodGet}, handler, opts...)
}

func (r *route) Post(handler faas.HttpMiddleware, opts ...MethodOption) {
	r.AddMethodHandler([]string{http.MethodPost}, handler, opts...)
}

func (r *route) Put(handler faas.HttpMiddleware, opts ...MethodOption) {
	r.AddMethodHandler([]string{http.MethodPut}, handler, opts...)
}

func (r *route) Patch(handler faas.HttpMiddleware, opts ...MethodOption) {
	r.AddMethodHandler([]string{http.MethodPatch}, handler, opts...)
}

func (r *route) Delete(handler faas.HttpMiddleware, opts ...MethodOption) {
	r.AddMethodHandler([]string{http.MethodDelete}, handler, opts...)
}

func (r *route) Options(handler faas.HttpMiddleware, opts ...MethodOption) {
	r.AddMethodHandler([]string{http.MethodOptions}, handler, opts...)
}

// Api Resource represents an HTTP API, capable of routing and securing incoming HTTP requests to handlers.
// path is the route path matcher e.g. '/home'. Supports path params via colon prefix e.g. '/customers/:customerId'
// handler the handler to register for callbacks.
//
// Note: to chain middleware use faas.ComposeHttpMiddlware()
type Api interface {
	Get(path string, handler faas.HttpMiddleware, opts ...MethodOption)
	Put(path string, handler faas.HttpMiddleware, opts ...MethodOption)
	Patch(path string, handler faas.HttpMiddleware, opts ...MethodOption)
	Post(path string, handler faas.HttpMiddleware, opts ...MethodOption)
	Delete(path string, handler faas.HttpMiddleware, opts ...MethodOption)
	Options(path string, handler faas.HttpMiddleware, opts ...MethodOption)
}

type api struct {
	name          string
	routes        map[string]Route
	m             *manager
	securityRules map[string]interface{}
	security      map[string][]string
}

func (m *manager) NewApi(name string, opts ...ApiOption) (Api, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	a := &api{
		name:   name,
		routes: map[string]Route{},
		m:      m,
	}

	// Apply options
	if opts != nil {
		for _, o := range opts {
			o(a)
		}
	}

	var secDefs map[string]*v1.ApiSecurityDefinition = nil
	var security map[string]*v1.ApiScopes = nil

	// Apply security rules
	if a.securityRules != nil {
		secDefs = make(map[string]*v1.ApiSecurityDefinition)
		for n, def := range a.securityRules {
			if jwt, ok := def.(JwtSecurityRule); ok {
				secDefs[n] = &v1.ApiSecurityDefinition{
					Definition: &v1.ApiSecurityDefinition_Jwt{
						Jwt: &v1.ApiSecurityDefinitionJwt{
							Issuer:    jwt.Issuer,
							Audiences: jwt.Audiences,
						},
					},
				}
			}
		}
	}

	// Apply security and scopes
	if a.security != nil {
		security = make(map[string]*v1.ApiScopes)
		for n, sec := range a.security {
			security[n] = &v1.ApiScopes{
				Scopes: sec,
			}
		}
	}

	// declare resource
	_, err = rsc.Declare(context.TODO(), &v1.ResourceDeclareRequest{
		Resource: &v1.Resource{
			Name: name,
			Type: v1.ResourceType_Api,
		},
		Config: &v1.ResourceDeclareRequest_Api{
			Api: &v1.ApiResource{
				SecurityDefinitions: secDefs,
				Security:            security,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}

// NewApi Registers a new API Resource.
//
// The returned API object can be used to register Routes and Methods, with Handlers.
func NewApi(name string, opts ...ApiOption) (Api, error) {
	return run.NewApi(name, opts...)
}

// Get adds a Get method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Get(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r, ok := a.routes[match]
	if !ok {
		r = a.m.NewRoute(a.name, match)
	}
	r.Get(handler, opts...)
	a.routes[match] = r
}

// Post adds a Post method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Post(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Post(handler, opts...)
	a.routes[match] = r
}

// Patch adds a Patch method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Patch(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Patch(handler, opts...)
	a.routes[match] = r
}

// Put adds a Put method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Put(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Put(handler, opts...)
	a.routes[match] = r
}

// Delete adds a Delete method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Delete(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Delete(handler, opts...)
	a.routes[match] = r
}

// Options adds an Options method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Options(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Options(handler, opts...)
	a.routes[match] = r
}
