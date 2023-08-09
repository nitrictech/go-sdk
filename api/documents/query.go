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

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

// Query - Query interface for Document Service
type Query interface {
	// Where - Append one or more expressions to the query
	Where(...*queryExpression) Query

	// Limit - limit the max result size of the query
	Limit(int) Query
	// FromPagingToken - Start from a given entry
	FromPagingToken(interface{}) Query

	// Fetch - Return paged values
	Fetch(ctx context.Context) (*FetchResult, error)

	// Stream - Return an iterator containing values
	Stream(ctx context.Context) (DocumentIter, error)
}

// Defacto Query interface implementation
type queryImpl struct {
	col            CollectionRef
	documentClient v1.DocumentServiceClient
	pagingToken    interface{}
	expressions    []*queryExpression
	limit          int
}

func (q *queryImpl) Where(qes ...*queryExpression) Query {
	q.expressions = append(q.expressions, qes...)

	return q
}

func (q *queryImpl) Limit(limit int) Query {
	q.limit = limit
	return q
}

func (q *queryImpl) FromPagingToken(token interface{}) Query {
	q.pagingToken = token
	return q
}

type FetchResult struct {
	Documents   []Document
	PagingToken interface{}
}

func (q *queryImpl) expressionsToWire() ([]*v1.Expression, error) {
	expressions := make([]*v1.Expression, 0, len(q.expressions))

	for _, e := range q.expressions {
		wexp, err := e.ToWire()
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, wexp)
	}

	return expressions, nil
}

func (q *queryImpl) Fetch(ctx context.Context) (*FetchResult, error) {
	// build the expressions list
	expressions, err := q.expressionsToWire()
	if err != nil {
		return nil, err
	}

	var token map[string]string = nil
	if q.pagingToken != nil {
		t, ok := q.pagingToken.(map[string]string)

		if !ok {
			return nil, errors.New(codes.InvalidArgument, "Query.Fetch: Paging Token invalid")
		}
		token = t
	}

	r, err := q.documentClient.Query(ctx, &v1.DocumentQueryRequest{
		Collection:  q.col.ToWire(),
		Expressions: expressions,
		Limit:       int32(q.limit),
		PagingToken: token,
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	docs := make([]Document, 0, len(r.GetDocuments()))

	for _, d := range r.GetDocuments() {
		ref, err := documentRefFromWireKey(q.documentClient, d.GetKey())
		if err != nil {
			// XXX: Potentially just log an error and continue
			return nil, err
		}

		docs = append(docs, &documentImpl{
			ref:     ref,
			content: d.Content.AsMap(),
		})
	}

	return &FetchResult{
		Documents:   docs,
		PagingToken: r.GetPagingToken(),
	}, nil
}

func (q *queryImpl) Stream(ctx context.Context) (DocumentIter, error) {
	// build the expressions list
	expressions, err := q.expressionsToWire()
	if err != nil {
		return nil, err
	}

	r, err := q.documentClient.QueryStream(ctx, &v1.DocumentQueryStreamRequest{
		Collection:  q.col.ToWire(),
		Expressions: expressions,
		Limit:       int32(q.limit),
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	// TODO: Return result iterator
	return &documentIterImpl{
		documentClient:       q.documentClient,
		documentStreamClient: r,
	}, nil
}

func newQuery(col CollectionRef, dc v1.DocumentServiceClient) Query {
	return &queryImpl{
		documentClient: dc,
		col:            col,
		expressions:    make([]*queryExpression, 0),
	}
}
