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

import "github.com/nitrictech/go-sdk/faas"

type Route interface {
	Get(handler ...faas.HttpMiddleware) error
	Patch(handler ...faas.HttpMiddleware) error
	Put(handler ...faas.HttpMiddleware) error
	Post(handler ...faas.HttpMiddleware) error
	Delete(handler ...faas.HttpMiddleware) error
}

type route struct {
	opts           faas.ApiWorkerOptions
	handlerBuilder faas.HandlerBuilder
}

func NewRoute(apiName, path string) Route {
	return &route{
		opts: faas.ApiWorkerOptions{
			ApiName:     apiName,
			Path:        path,
			HttpMethods: []string{},
		},
		handlerBuilder: faas.New(),
	}
}

func (r *route) optsWithMethod(method string) faas.ApiWorkerOptions {
	r.opts.HttpMethods = append(r.opts.HttpMethods, method)
	return r.opts
}

func (r *route) Get(handlers ...faas.HttpMiddleware) error {
	return r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("GET")).Start()
}

func (r *route) Post(handlers ...faas.HttpMiddleware) error {
	return r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("POST")).Start()
}

func (r *route) Put(handlers ...faas.HttpMiddleware) error {
	return r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("PUT")).Start()
}

func (r *route) Patch(handlers ...faas.HttpMiddleware) error {
	return r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("PATCH")).Start()
}

func (r *route) Delete(handlers ...faas.HttpMiddleware) error {
	return r.handlerBuilder.Http(handlers...).WithApiWorkerOpts(r.optsWithMethod("DELETE")).Start()
}

//	mainApi := nitric.Api("main")
//
//	mainApi.Get("/hello/:name", func(ctx *faas.HttpContext, next *faas.HttpHandler) *faas.HttpContext {
//	  // implement
//
//    return next(ctx)
//	})

type Api interface {
	Get(string, ...faas.HttpMiddleware) error
	Put(string, ...faas.HttpMiddleware) error
	Patch(string, ...faas.HttpMiddleware) error
	Post(string, ...faas.HttpMiddleware) error
	Delete(string, ...faas.HttpMiddleware) error
}

type api struct {
	name   string
	routes map[string]Route
}

func NewApi(name string) Api {
	return &api{
		name:   name,
		routes: map[string]Route{},
	}
}

func (a *api) Get(match string, handlers ...faas.HttpMiddleware) error {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	return r.Get(handlers...)
}

func (a *api) Post(match string, handlers ...faas.HttpMiddleware) error {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	return r.Post(handlers...)
}

func (a *api) Patch(match string, handlers ...faas.HttpMiddleware) error {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	return r.Patch(handlers...)
}

func (a *api) Put(match string, handlers ...faas.HttpMiddleware) error {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	return r.Put(handlers...)
}

func (a *api) Delete(match string, handlers ...faas.HttpMiddleware) error {
	r, ok := a.routes[match]
	if !ok {
		r = NewRoute(a.name, match)
	}
	return r.Delete(handlers...)
}
