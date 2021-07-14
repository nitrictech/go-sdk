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

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
)

// Query - Query interface for Document Service
type Query interface {
	// Where - Append an expression to the query
	Where(...*queryExpression) Query

	// Limit - limit the max result size of the query
	Limit(int) Query
	// FromPagingToken - Start from a given entry
	FromPagingToken(interface{}) Query

	// Fetch - Return paged values
	Fetch() (*FetchResult, error)

	// Stream - Return an iterator containing values
	Stream() (DocumentIter, error)
}

// Defacto Query interface implementation
type queryImpl struct {
	col   CollectionRef
	dc    v1.DocumentServiceClient
	pt    interface{}
	exps  []*queryExpression
	limit int
}

func (q *queryImpl) Where(qes ...*queryExpression) Query {
	q.exps = append(q.exps, qes...)

	return q
}

func (q *queryImpl) Limit(limit int) Query {
	q.limit = limit
	return q
}

func (q *queryImpl) FromPagingToken(token interface{}) Query {
	q.pt = token
	return q
}

type FetchResult struct {
	Documents   []Document
	PagingToken interface{}
}

func (q *queryImpl) Fetch() (*FetchResult, error) {
	// build the expressions list
	expressions := make([]*v1.Expression, 0, len(q.exps))

	for _, e := range q.exps {
		wexp, err := e.toWire()

		if err != nil {
			return nil, err
		}

		expressions = append(expressions, wexp)
	}

	var token map[string]string = nil
	if q.pt != nil {
		t, ok := q.pt.(map[string]string)

		if !ok {
			return nil, fmt.Errorf("Invalid paging token provided!")
		}
		token = t
	}

	r, err := q.dc.Query(context.TODO(), &v1.DocumentQueryRequest{
		Collection:  q.col.toWire(),
		Expressions: expressions,
		Limit:       int32(q.limit),
		PagingToken: token,
	})

	if err != nil {
		return nil, err
	}

	docs := make([]Document, 0, len(r.GetDocuments()))

	for _, d := range r.GetDocuments() {
		docs = append(docs, &documentImpl{
			content: d.Content.AsMap(),
		})
	}

	return &FetchResult{
		Documents:   docs,
		PagingToken: r.GetPagingToken(),
	}, nil
}

func (q *queryImpl) Stream() (DocumentIter, error) {
	// build the expressions list
	expressions := make([]*v1.Expression, 0, len(q.exps))

	for _, e := range q.exps {
		wexp, err := e.toWire()

		if err != nil {
			return nil, err
		}

		expressions = append(expressions, wexp)
	}

	r, err := q.dc.QueryStream(context.TODO(), &v1.DocumentQueryStreamRequest{
		Collection:  q.col.toWire(),
		Expressions: expressions,
		Limit:       int32(q.limit),
	})

	if err != nil {
		return nil, err
	}

	// TODO: Return result iterator
	return &documentIterImpl{
		str: r,
	}, nil
}

func newQuery(col CollectionRef, dc v1.DocumentServiceClient) Query {
	return &queryImpl{
		dc:   dc,
		col:  col,
		exps: make([]*queryExpression, 0),
	}
}
