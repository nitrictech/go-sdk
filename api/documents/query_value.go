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

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
)

type value struct {
	number_value *int
	string_value *string
	double_value *float64
	bool_value   *bool
}

// NumberValue - Create a NumberValue for a Query
func NumberValue(val int) *value {
	return &value{
		number_value: &val,
	}
}

// StringValue - Create a StringValue for a Query
func StringValue(val string) *value {
	return &value{
		string_value: &val,
	}
}

// DoubleValue - Create a DoubleValue for a Query
func DoubleValue(val float64) *value {
	return &value{
		double_value: &val,
	}
}

func BoolValue(val bool) *value {
	return &value{
		bool_value: &val,
	}
}

func (v *value) toWire() (*v1.ExpressionValue, error) {
	if v.number_value != nil {
		return &v1.ExpressionValue{
			Kind: &v1.ExpressionValue_IntValue{
				IntValue: int64(*v.number_value),
			},
		}, nil
	} else if v.double_value != nil {
		return &v1.ExpressionValue{
			Kind: &v1.ExpressionValue_DoubleValue{
				DoubleValue: *v.double_value,
			},
		}, nil
	} else if v.string_value != nil {
		return &v1.ExpressionValue{
			Kind: &v1.ExpressionValue_StringValue{
				StringValue: *v.string_value,
			},
		}, nil
	} else if v.bool_value != nil {
		return &v1.ExpressionValue{
			Kind: &v1.ExpressionValue_BoolValue{
				BoolValue: *v.bool_value,
			},
		}, nil
	}

	return nil, fmt.Errorf("Invalid query value")
}
