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
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// DocumentRef - Represents a reference to a document
type DocumentRef interface {
	// toWireKey - Translate this document ready for on-wire transport
	toWireKey() *v1.Key

	// Get - Retrieve the value of the document
	Get() (Document, error)

	// Set - Sets the value of the document
	Set(map[string]interface{}) error

	// Delete - Deletes the document
	Delete() error

	// Collection - Retrieve a child collection of this document
	Collection(string) (CollectionRef, error)
}

type documentRefImpl struct {
	// A shared reference to the top level document service client
	dc v1.DocumentServiceClient
	// A reference to this documents collection
	col CollectionRef
	// The key for this document
	key string
}

func (d *documentRefImpl) toWireKey() *v1.Key {
	return &v1.Key{
		Id:         d.key,
		Collection: d.col.toWire(),
	}
}

func (d *documentRefImpl) Collection(c string) (CollectionRef, error) {
	if d.col.Parent() != nil {
		return nil, fmt.Errorf("Nested sub-collections are currently not supported")
	}

	return &collectionRefImpl{
		name:           c,
		dc:             d.dc,
		parentDocument: d,
	}, nil
}

func (d *documentRefImpl) Delete() error {
	_, err := d.dc.Delete(context.TODO(), &v1.DocumentDeleteRequest{
		Key: &v1.Key{
			Collection: d.col.toWire(),
			Id:         d.key,
		},
	})

	return err
}

func (d *documentRefImpl) Set(content map[string]interface{}) error {
	sv, err := structpb.NewStruct(content)

	fmt.Println("struct", sv)

	if err != nil {
		return err
	}

	_, err = d.dc.Set(context.TODO(), &v1.DocumentSetRequest{
		Key: &v1.Key{
			Collection: d.col.toWire(),
			Id:         d.key,
		},
		Content: sv,
	})

	return err
}

type DecodeOption interface {
	Apply(c *mapstructure.DecoderConfig)
}

// Get -
func (d *documentRefImpl) Get() (Document, error) {
	res, err := d.dc.Get(context.TODO(), &v1.DocumentGetRequest{
		Key: &v1.Key{
			Collection: d.col.toWire(),
			Id:         d.key,
		},
	})

	if err != nil {
		return nil, err
	}

	return &documentImpl{
		content: res.GetDocument().GetContent().AsMap(),
	}, nil
}
