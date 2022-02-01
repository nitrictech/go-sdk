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
}

type route struct {
	opts           faas.ApiWorkerOptions
	handlerBuilder faas.HandlerBuilder
}

func NewRoute(apiName, apiPath string) Route {
	return run.NewRoute(apiName, apiPath)
}

func (m *manager) NewRoute(apiName, apiPath string) Route {
	f := faas.New()
	m.addStarter("route:"+path.Join(apiName, apiPath), f)

	return &route{
		opts: faas.ApiWorkerOptions{
			ApiName:     apiName,
			Path:        apiPath,
			HttpMethods: []string{},
		},
		handlerBuilder: f,
	}
}

func (r *route) optsWithMethod(method string) faas.ApiWorkerOptions {
	r.opts.HttpMethods = append(r.opts.HttpMethods, method)
	return r.opts
}

func (r *route) Get(handlers ...faas.HttpMiddleware) {
	r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("GET"))
}

func (r *route) Post(handlers ...faas.HttpMiddleware) {
	r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("POST"))
}

func (r *route) Put(handlers ...faas.HttpMiddleware) {
	r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("PUT"))
}

func (r *route) Patch(handlers ...faas.HttpMiddleware) {
	r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("PATCH"))
}

func (r *route) Delete(handlers ...faas.HttpMiddleware) {
	r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("DELETE"))
}

//	mainApi := nitric.Api("main")
//
//	mainApi.Get("/hello/:name", func(ctx *faas.HttpContext, next *faas.HttpHandler) *faas.HttpContext {
//	  // implement
//
//    return next(ctx)
//	})

type Api interface {
	Get(string, ...faas.HttpMiddleware)
	Put(string, ...faas.HttpMiddleware)
	Patch(string, ...faas.HttpMiddleware)
	Post(string, ...faas.HttpMiddleware)
	Delete(string, ...faas.HttpMiddleware)
}

type api struct {
	name   string
	routes map[string]Route
}

func (m *manager) NewApi(name string) Api {
	return &api{
		name:   name,
		routes: map[string]Route{},
	}
}

func NewApi(name string) Api {
	return run.NewApi(name)
}

func (a *api) Get(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Get(handlers...)
}

func (a *api) Post(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Post(handlers...)
}

func (a *api) Patch(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Patch(handlers...)
}

func (a *api) Put(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Put(handlers...)
}

func (a *api) Delete(match string, handlers ...faas.HttpMiddleware) {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	r.Delete(handlers...)
}
