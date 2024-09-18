package http

type Response struct {
	Status  int
	Headers map[string][]string
	Body    []byte
}
