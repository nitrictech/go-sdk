package faas

import "net/http"

// NitricResponse - represents the results of calling a function.
type NitricResponse struct {
	Headers map[string]string
	Status  int
	Body    []byte
}

// writeHttpResponse - writes a HTTP response from a NitricResponse
func (n *NitricResponse) writeHTTPResponse(w http.ResponseWriter) {
	for k, v := range n.Headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(n.Status)
	w.Write(n.Body)
}
