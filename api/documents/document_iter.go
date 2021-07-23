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

// DocumentIter - An iterator for lazy document retrieval
type DocumentIter interface {
	// Next - Retrieve the next document in the iterator
	Next() (Document, error)
}

type documentIterImpl struct {
	dc  v1.DocumentServiceClient
	str v1.DocumentService_QueryStreamClient
}

func (i *documentIterImpl) Next() (Document, error) {
	res, err := i.str.Recv()

	if err != nil {
		return nil, err
	}

	ref, err := documentRefFromWireKey(i.dc, res.GetDocument().GetKey())
	if err != nil {
		return nil, err
	}

	return &documentImpl{
		ref:     ref,
		content: res.GetDocument().GetContent().AsMap(),
	}, nil
}
