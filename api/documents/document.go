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

import "github.com/mitchellh/mapstructure"

// Document - Interface for interacing with document API results
type Document interface {
	// Ref - Retrieve the original DocumentRef for this document
	Ref() DocumentRef
	// Content - Retrieve the documents content as a Map
	Content() map[string]interface{}
	// Decode - Decode document content into the given reference
	Decode(interface{}, ...DecodeOption) error
}

type documentImpl struct {
	ref     DocumentRef
	content map[string]interface{}
}

func (d *documentImpl) Ref() DocumentRef {
	return d.ref
}

func (d *documentImpl) Content() map[string]interface{} {
	return d.content
}

// Decode -
func (d *documentImpl) Decode(val interface{}, opts ...DecodeOption) error {
	decoderConfig := mapstructure.DecoderConfig{
		//DecodeHook:       nil,
		ErrorUnused: true, // Default behavior is to error when keys are missing from the output interface{}
		//ZeroFields:       false,
		//WeaklyTypedInput: false,
		//Squash:           false,
		//Metadata:         nil,
		Result: val,
		//TagName:          "",
	}

	// Apply additional options
	for _, opt := range opts {
		opt.Apply(&decoderConfig)
	}

	// Decode the value into the object
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return err
	}
	return decoder.Decode(d.content)
}
