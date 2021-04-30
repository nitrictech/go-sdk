// Copyright 2021 Nitric Pty Ltd.
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

package faas

// SourceType - enum of the possible sources for a Nitric request.
type SourceType string

const (
	// Request - HTTP Request Source Type
	Request SourceType = "REQUEST"
	// Subscription - Topic Subscription Source Type
	Subscription = "SUBSCRIPTION"
	// Unknown - Unknown Source Types, used when the source can't be determined.
	Unknown = "UNKNOWN"
)

// Each of the source type consts above must be included here.
var sourceTypes []SourceType = []SourceType{
	Request,
	Subscription,
	Unknown,
}

// sourceTypeFromString - converts a string, typically the x-nitric-source-type header, into a SourceType
func sourceTypeFromString(s string) SourceType {
	for _, t := range sourceTypes {
		if s == string(t) {
			return t
		}
	}
	// Default to Unknown if the source type isn't one that has been defined.
	return Unknown
}

func (p SourceType) String() string {
	x := p
	for _, v := range sourceTypes {
		if v == x {
			return string(x)
		}
	}
	return Unknown // This will only happen if manually changed.
}

// NitricContext - Represents the contextual metadata for a Nitric function request.
type NitricContext struct {
	requestID   string
	source      string
	sourceType  SourceType
	payloadType string
}

// GetRequestID - return the request id of the request.
func (c *NitricContext) GetRequestID() string {
	return c.requestID
}

// GetSource - return the source of the request.
func (c *NitricContext) GetSource() string {
	return c.source
}

// GetSourceType - return the source type of the request
func (c *NitricContext) GetSourceType() SourceType {
	return c.sourceType
}

// GetPayloadType - return the payload type of the request payload. Typically a typehint.
func (c *NitricContext) GetPayloadType() string {
	return c.payloadType
}
