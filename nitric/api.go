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
	"net/http"
	"path"
	"strings"

	v1 "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
)

// Route providers convenience functions to register a handler in a single method.
type Route interface {
	All(handler faas.HttpMiddleware, opts ...MethodOption)
	Get(handler faas.HttpMiddleware, opts ...MethodOption)
	Patch(handler faas.HttpMiddleware, opts ...MethodOption)
	Put(handler faas.HttpMiddleware, opts ...MethodOption)
	Post(handler faas.HttpMiddleware, opts ...MethodOption)
	Delete(handler faas.HttpMiddleware, opts ...MethodOption)
	Options(handler faas.HttpMiddleware, opts ...MethodOption)
}

type route struct {
	path       string
	apiName    string
	middleware faas.HttpMiddleware
	manager    *manager
}

func composeRouteMiddleware(apiMiddleware faas.HttpMiddleware, routeMiddleware []faas.HttpMiddleware) faas.HttpMiddleware {
	if apiMiddleware != nil && len(routeMiddleware) > 0 {
		return faas.ComposeHttpMiddleware(apiMiddleware, faas.ComposeHttpMiddleware(routeMiddleware...))
	}

	if len(routeMiddleware) > 0 {
		return faas.ComposeHttpMiddleware(routeMiddleware...)
	}

	return apiMiddleware
}

func (a *api) NewRoute(match string, middleware ...faas.HttpMiddleware) Route {
	r, ok := a.routes[match]
	if !ok {
		r = &route{
			manager:    a.manager,
			path:       path.Join(a.path, match),
			apiName:    a.name,
			middleware: composeRouteMiddleware(a.middleware, middleware),
		}
	}

	return r
}

func (r *route) AddMethodHandler(methods []string, handler faas.HttpMiddleware, opts ...MethodOption) {
	bName := path.Join(r.apiName, r.path, strings.Join(methods, "-"))

	mo := &methodOptions{}
	for _, o := range opts {
		o(mo)
	}

	b := r.manager.getBuilder(bName)
	if b == nil {
		b = faas.New().WithApiWorkerOpts(faas.ApiWorkerOptions{
			ApiName:          r.apiName,
			Path:             r.path,
			Security:         mo.security,
			SecurityDisabled: mo.securityDisabled,
		})
	}

	composedHandler := handler
	if r.middleware != nil {
		composedHandler = faas.ComposeHttpMiddleware(r.middleware, handler)
	}

	for _, m := range methods {
		b.Http(m, composedHandler)
	}

	r.manager.addBuilder(bName, b)
	r.manager.addWorker("route:"+bName, b)
}

func (r *route) All(handler faas.HttpMiddleware, opts ...MethodOption) {
	r.AddMethodHandler([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}, handler, opts...)
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
	NewRoute(path string, middleware ...faas.HttpMiddleware) Route
	Details(ctx context.Context) (*ApiDetails, error)
}

type ApiDetails struct {
	Details
	URL string
}

type api struct {
	name          string
	routes        map[string]Route
	manager       *manager
	securityRules map[string]interface{}
	security      map[string][]string
	path          string
	middleware    faas.HttpMiddleware
}

func (m *manager) newApi(name string, opts ...ApiOption) (Api, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	a := &api{
		name:    name,
		routes:  map[string]Route{},
		manager: m,
	}

	// Apply options
	for _, o := range opts {
		o(a)
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
	return defaultManager.newApi(name, opts...)
}

func (a *api) Details(ctx context.Context) (*ApiDetails, error) {
	rsc, err := a.manager.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	resp, err := rsc.Details(ctx, &v1.ResourceDetailsRequest{
		Resource: &v1.Resource{
			Type: v1.ResourceType_Api,
			Name: a.name,
		},
	})
	if err != nil {
		return nil, err
	}

	d := &ApiDetails{
		Details: Details{
			ID:       resp.Id,
			Provider: resp.Provider,
			Service:  resp.Service,
		},
	}
	if resp.GetApi() != nil {
		d.URL = resp.GetApi().GetUrl()
	}

	return d, nil
}

// Get adds a Get method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Get(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r := a.NewRoute(match)

	r.Get(handler, opts...)
	a.routes[match] = r
}

// Post adds a Post method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Post(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r := a.NewRoute(match)

	r.Post(handler, opts...)
	a.routes[match] = r
}

// Patch adds a Patch method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Patch(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r := a.NewRoute(match)

	r.Patch(handler, opts...)
	a.routes[match] = r
}

// Put adds a Put method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Put(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r := a.NewRoute(match)

	r.Put(handler, opts...)
	a.routes[match] = r
}

// Delete adds a Delete method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Delete(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r := a.NewRoute(match)

	r.Delete(handler, opts...)
	a.routes[match] = r
}

// Options adds an Options method handler to the path with any specified opts.
// Note: to chain middleware use faas.ComposeHttpMiddlware()
func (a *api) Options(match string, handler faas.HttpMiddleware, opts ...MethodOption) {
	r := a.NewRoute(match)

	r.Options(handler, opts...)
	a.routes[match] = r
}
