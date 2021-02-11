package faas

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// NitricRequest - represents a request to trigger a function, with payload and context required to execute that function.
type NitricRequest struct {
	context NitricContext
	payload []byte
}

// GetContext - return the context of a request, with metadata about that request.
func (n *NitricRequest) GetContext() NitricContext {
	return n.context
}

// GetPayload - return the []byte payload of the request.
func (n *NitricRequest) GetPayload() []byte {
	return n.payload
}

// GetStruct - Unmarshals the request body from JSON to the provided interface{}
func (n *NitricRequest) GetStruct(object interface{}) error {
	return json.Unmarshal(n.payload, object)
}

// contextFromHeaders - converts standard nitric HTTP headers into a context struct.
func contextFromHeaders(h http.Header) NitricContext {
	return NitricContext{
		requestID:   h.Get("x-nitric-request-id"),
		sourceType:  sourceTypeFromString(h.Get("x-nitric-source-type")),
		source:      h.Get("x-nitric-source"),
		payloadType: h.Get("x-nitric-payload-type"),
	}
}

// fromHttpRequest - converts a standard nitric HTTP request into a NitricRequest to be passed to functions.
func fromHttpRequest(r *http.Request) (*NitricRequest, error) {
	context := contextFromHeaders(r.Header)

	payload, err := ioutil.ReadAll(r.Body)

	if err == nil {
		return &NitricRequest{
			context: context,
			payload: payload,
		}, nil
	}

	return nil, err
}
