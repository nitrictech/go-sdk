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
	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// DocumentRef - Represents a reference to a document
type DocumentRef interface {
	// toWireKey - Translate this document ready for on-wire transport
	toWireKey() *v1.Key

	Parent() CollectionRef

	Id() string

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
	// The id for this document
	id string
}

// Construct a document reference from the wire
func documentRefFromWireKey(dc v1.DocumentServiceClient, k *v1.Key) (DocumentRef, error) {
	if dc != nil && k != nil {
		col, err := collectionRefFromWire(dc, k.GetCollection())

		if err != nil {
			return nil, err
		}

		return &documentRefImpl{
			dc:  dc,
			col: col,
			id:  k.GetId(),
		}, nil
	} else {
		if dc == nil {
			return nil, errors.New(
				codes.Internal,
				"documentRefFromWireKey: provide non-nil DocumentServiceClient",
			)
		}

		return nil, errors.New(
			codes.Internal,
			"documentRefFromWireKey: provide non-nil Key",
		)
	}
}

func (d *documentRefImpl) toWireKey() *v1.Key {
	return &v1.Key{
		Id:         d.id,
		Collection: d.col.toWire(),
	}
}

// Id - The documents ID
func (d *documentRefImpl) Id() string {
	return d.id
}

// Parent - Gets the Parent collection reference for this document
func (d *documentRefImpl) Parent() CollectionRef {
	return d.col
}

// Collection - Gets a subcollection for this document
func (d *documentRefImpl) Collection(c string) (CollectionRef, error) {
	if d.col.Parent() != nil {
		return nil, errors.New(
			codes.InvalidArgument,
			fmt.Sprintf("DocumentRef.Collection: Maximum collection depth: %d exceeded", MaxCollectionDepth),
		)
	}

	return &collectionRefImpl{
		name:           c,
		dc:             d.dc,
		parentDocument: d,
	}, nil
}

// Delete - Deletes the document this reference refers to if it exists
func (d *documentRefImpl) Delete() error {
	_, err := d.dc.Delete(context.TODO(), &v1.DocumentDeleteRequest{
		Key: d.toWireKey(),
	})

	return err
}

// Set - Sets the contents of the document this reference refers to
func (d *documentRefImpl) Set(content map[string]interface{}) error {
	sv, err := structpb.NewStruct(content)

	if err != nil {
		return errors.NewWithCause(
			codes.Internal,
			"DocumentRef.Set: Unable to create protobuf struct",
			err,
		)
	}

	if _, err = d.dc.Set(context.TODO(), &v1.DocumentSetRequest{
		Key:     d.toWireKey(),
		Content: sv,
	}); err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

type DecodeOption interface {
	Apply(c *mapstructure.DecoderConfig)
}

// Get - Retrieves the Document this reference refers to if it exists
func (d *documentRefImpl) Get() (Document, error) {
	res, err := d.dc.Get(context.TODO(), &v1.DocumentGetRequest{
		Key: d.toWireKey(),
	})

	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	return &documentImpl{
		ref:     d,
		content: res.GetDocument().GetContent().AsMap(),
	}, nil
}
