package http

type Request interface {
	Method() string
	Path() string
	Data() []byte
	Query() map[string][]string
	Headers() map[string][]string
	PathParams() map[string]string
}

type RequestImpl struct {
	method     string
	path       string
	data       []byte
	query      map[string][]string
	headers    map[string][]string
	pathParams map[string]string
}

func (h *RequestImpl) Method() string {
	return h.method
}

func (h *RequestImpl) Path() string {
	return h.path
}

func (h *RequestImpl) Data() []byte {
	return h.data
}

func (h *RequestImpl) Query() map[string][]string {
	return h.query
}

func (h *RequestImpl) Headers() map[string][]string {
	return h.headers
}

func (h *RequestImpl) PathParams() map[string]string {
	return h.pathParams
}
