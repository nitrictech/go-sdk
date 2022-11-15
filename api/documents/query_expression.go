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
	"fmt"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	v1 "github.com/nitrictech/go-sdk/nitric/v1"
)

// QueryOp - Enum for query operations
type queryOp string

const (
	queryOp_EQ         queryOp = "=="
	queryOp_GT         queryOp = ">"
	queryOp_GE         queryOp = ">="
	queryOp_LT         queryOp = "<"
	queryOp_LE         queryOp = "<="
	queryOp_StartsWith queryOp = "startsWith"
)

// IsValid - Determine if the QueryOp is a valid value
func (o queryOp) IsValid() error {
	switch o {
	case queryOp_EQ, queryOp_GT, queryOp_GE, queryOp_LT, queryOp_LE, queryOp_StartsWith:
		return nil
	default:
		return errors.New(
			codes.InvalidArgument,
			fmt.Sprintf("queryOp.IsValid: invalid query operation (%s)", o),
		)
	}
}

type queryExpression struct {
	field string
	op    queryOp
	val   *value
}

func (q *queryExpression) ToWire() (*v1.Expression, error) {
	if q.field == "" {
		return nil, errors.New(
			codes.InvalidArgument,
			"queryExpress.ToWire: provide non-blank field name",
		)
	}

	if err := q.op.IsValid(); err != nil {
		return nil, err
	}

	wv, err := q.val.ToWire()
	if err != nil {
		return nil, err
	}

	return &v1.Expression{
		Operand:  q.field,
		Operator: string(q.op),
		Value:    wv,
	}, nil
}

type queryExpressionBuilder struct {
	field string
}

func (q *queryExpressionBuilder) Eq(val *value) *queryExpression {
	return &queryExpression{
		field: q.field,
		op:    queryOp_EQ,
		val:   val,
	}
}

func (q *queryExpressionBuilder) Lt(val *value) *queryExpression {
	return &queryExpression{
		field: q.field,
		op:    queryOp_LT,
		val:   val,
	}
}

func (q *queryExpressionBuilder) Le(val *value) *queryExpression {
	return &queryExpression{
		field: q.field,
		op:    queryOp_LE,
		val:   val,
	}
}

func (q *queryExpressionBuilder) Gt(val *value) *queryExpression {
	return &queryExpression{
		field: q.field,
		op:    queryOp_GT,
		val:   val,
	}
}

func (q *queryExpressionBuilder) Ge(val *value) *queryExpression {
	return &queryExpression{
		field: q.field,
		op:    queryOp_GE,
		val:   val,
	}
}

func (q *queryExpressionBuilder) StartsWith(val *value) *queryExpression {
	return &queryExpression{
		field: q.field,
		op:    queryOp_StartsWith,
		val:   val,
	}
}

func Condition(field string) *queryExpressionBuilder {
	return &queryExpressionBuilder{
		field: field,
	}
}
