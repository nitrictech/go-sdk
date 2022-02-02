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

import (
	v1 "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
)

// Collection
type CollectionRef interface {
	Name() string
	Doc(string) DocumentRef
	Query() Query
	Parent() DocumentRef
	Collection(string) CollectionGroupRef
	ToWire() *v1.Collection
}

type collectionRefImpl struct {
	name           string
	dc             v1.DocumentServiceClient
	parentDocument DocumentRef
}

// Query - Returns a Query builder to construct queries against this collection
func (c *collectionRefImpl) Query() Query {
	return newQuery(c, c.dc)
}

// Name - Returns the name of the collection
func (c *collectionRefImpl) Name() string {
	return c.name
}

// Doc - Return a document reference for this collection
func (c *collectionRefImpl) Doc(key string) DocumentRef {
	return &documentRefImpl{
		id:  key,
		dc:  c.dc,
		col: c,
	}
}

// Parent - Retrieve the parent document of this collection
func (c *collectionRefImpl) Parent() DocumentRef {
	return c.parentDocument
}

// Collection - Creates a collection group reference from a collection
func (c *collectionRefImpl) Collection(name string) CollectionGroupRef {
	return &collectionGroupRefImpl{
		parent: fromColRef(c, c.dc),
		dc:     c.dc,
		name:   name,
	}
}

// ToWire - tranlates a Collection for on-wire transport
func (c *collectionRefImpl) ToWire() *v1.Collection {
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

// converts a wire collection to a collection reference
func collectionRefFromWire(dc v1.DocumentServiceClient, c *v1.Collection) (CollectionRef, error) {
	if dc == nil {
		return nil, errors.New(codes.Internal, "collectionRefFromWire: missing DocumentServiceClient")
	}

	if c.GetParent() == nil {
		return &collectionRefImpl{
			name: c.GetName(),
			dc:   dc,
		}, nil
	} else {
		pd, err := documentRefFromWireKey(dc, c.GetParent())
		if err != nil {
			return nil, errors.NewWithCause(codes.Internal, "collectionRefFromWire", err)
		}

		return &collectionRefImpl{
			name:           c.GetName(),
			dc:             dc,
			parentDocument: pd,
		}, nil
	}
}
