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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/go-sdk/faas"
)

var _ = Describe("api", func() {
	Context("New", func() {
		It("can register multiple methods on routes", func() {
			m := &manager{
				blockers: map[string]Starter{},
				builders: map[string]faas.HandlerBuilder{},
			}
			a := &api{
				name:   "testApi",
				routes: map[string]Route{},
				m:      m,
			}
			a.Get("objects/", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Post("objects/", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Get("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Put("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Patch("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Delete("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })

			Expect(m.blockers["route:testApi/objects/:id"]).ToNot(BeNil())
			Expect(m.blockers["route:testApi/objects"]).ToNot(BeNil())
			Expect(m.builders["testApi/objects/:id"].String()).To(Equal("Api:testApi, path:objects/:id methods:[DELETE,GET,PATCH,PUT]"))
			Expect(m.builders["testApi/objects"].String()).To(Equal("Api:testApi, path:objects/ methods:[GET,POST]"))
		})
	})
})
