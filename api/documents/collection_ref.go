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

package documents

import v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"

// Collection
type CollectionRef interface {
	Doc(string) DocumentRef
	Query() Query
	Parent() DocumentRef
	toWire() *v1.Collection
}

type collectionRefImpl struct {
	name           string
	dc             v1.DocumentServiceClient
	parentDocument DocumentRef
}

func (c *collectionRefImpl) Query() Query {
	return newQuery(c, c.dc)
}

// Doc - Return a document reference for this collection
func (c *collectionRefImpl) Doc(key string) DocumentRef {
	return &documentRefImpl{
		key: key,
		dc:  c.dc,
		col: c,
	}
}

func (c *collectionRefImpl) Parent() DocumentRef {
	return c.parentDocument
}

// toWire - tranlates a Collection for on-wire transport
func (c *collectionRefImpl) toWire() *v1.Collection {
	if c.parentDocument != nil {
		return &v1.Collection{
			Name:   c.name,
			Parent: c.parentDocument.toWireKey(),
		}
	} else {
		return &v1.Collection{
			Name: c.name,
		}
	}
}
