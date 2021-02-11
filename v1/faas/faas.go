package faas

import "net/http"

// NitricFunction - a function built using Nitric, to be executed
type NitricFunction func(*NitricRequest) *NitricResponse

// Start - Starts accepting requests for the provided NitricFunction
//
// This should be the only method called in the 'main' method of your entrypoint package
func Start(f NitricFunction) {
	// Listen on the perscribed Nitric Application port (from env variables)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Translate the given HTTP request to a NitricRequest

		nr, err := fromHttpRequest(r)

		if err != nil {
			// Return a HTTP error here
			// do not call the inner function
			w.Header().Add("ContentType", "text/plain")
			w.WriteHeader(400)
			w.Write([]byte("Unable to read provided payload"))
			return
		}

		response := f(nr)

		// Write the reponse
		response.writeHTTPResponse(w)
	})
}
