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
	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

// Documents
type Documents interface {
	// Collection - Create a new collection reference
	Collection(string) CollectionRef
}

type documentsImpl struct {
	dc v1.DocumentServiceClient
}

// Collection - Create a new collection reference
func (d *documentsImpl) Collection(col string) CollectionRef {
	return &collectionRefImpl{
		dc:             d.dc,
		name:           col,
		parentDocument: nil,
	}
}

func New() (Documents, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)

	if err != nil {
		return nil, errors.NewWithCause(codes.Internal, "documents.New", err)
	}

	dc := v1.NewDocumentServiceClient(conn)

	return &documentsImpl{
		dc: dc,
	}, nil
}
