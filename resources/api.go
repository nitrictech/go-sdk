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
	"path"

	"github.com/nitrictech/go-sdk/faas"
)

type Route interface {
	Get(handler ...faas.HttpMiddleware)
	Patch(handler ...faas.HttpMiddleware)
	Put(handler ...faas.HttpMiddleware)
	Post(handler ...faas.HttpMiddleware)
	Delete(handler ...faas.HttpMiddleware)
	Options(handler ...faas.HttpMiddleware)
}

type route struct {
	builderName string
	m           *manager
}

func NewRoute(apiName, apiPath string) Route {
	return run.NewRoute(apiName, apiPath)
}

func (m *manager) NewRoute(apiName, apiPath string) Route {
	rName := path.Join(apiName, apiPath)
	_, ok := m.builders[rName]
	if !ok {
		m.builders[rName] = faas.New().WithApiWorkerOpts(faas.ApiWorkerOptions{
			ApiName: apiName,
			Path:    apiPath,
		})
	}

	return &route{
		m:           m,
		builderName: rName,
	}
}

func (r *route) addMethodHandler(method string, handlers ...faas.HttpMiddleware) {
	b := r.m.builders[r.builderName]
	b.Http(method, handlers...)
	r.m.addStarter("route:"+r.builderName, b)
	r.m.builders[r.builderName] = b
}

func (r *route) Get(handlers ...faas.HttpMiddleware) {
	r.addMethodHandler("GET", handlers...)
}

func (r *route) Post(handlers ...faas.HttpMiddleware) {
	r.addMethodHandler("POST", handlers...)
}

func (r *route) Put(handlers ...faas.HttpMiddleware) {
	r.addMethodHandler("PUT", handlers...)
}

func (r *route) Patch(handlers ...faas.HttpMiddleware) {
	r.addMethodHandler("PATCH", handlers...)
}

func (r *route) Delete(handlers ...faas.HttpMiddleware) {
	r.addMethodHandler("DELETE", handlers...)
}

func (r *route) Options(handlers ...faas.HttpMiddleware) {
	r.addMethodHandler("OPTIONS", handlers...)
}

type Api interface {
	Get(string, ...faas.HttpMiddleware)
	Put(string, ...faas.HttpMiddleware)
	Patch(string, ...faas.HttpMiddleware)
	Post(string, ...faas.HttpMiddleware)
	Delete(string, ...faas.HttpMiddleware)
	Options(string, ...faas.HttpMiddleware)
}

type api struct {
	name   string
	routes map[string]Route
	m      *manager
}

func (m *manager) NewApi(name string) Api {
	return &api{
		name:   name,
		routes: map[string]Route{},
		m:      m,
	}
}

func NewApi(name string) Api {
	return run.NewApi(name)
}

func (a *api) Get(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = a.m.NewRoute(a.name, match)
	}
	r.Get(handlers...)
	a.routes[match] = r
}

func (a *api) Options(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = a.m.NewRoute(a.name, match)
	}
	r.Options(handlers...)
	a.routes[match] = r
}

func (a *api) Post(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Post(handlers...)
	a.routes[match] = r
}

func (a *api) Patch(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Patch(handlers...)
	a.routes[match] = r
}

func (a *api) Put(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Put(handlers...)
	a.routes[match] = r
}

func (a *api) Delete(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Delete(handlers...)
	a.routes[match] = r
}
