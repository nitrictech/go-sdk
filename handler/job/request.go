package job

type Request interface {
	JobName() string
	Data() map[string]interface{}
}

type request struct {
	jobName string
	data    map[string]interface{}
}

func (m *request) JobName() string {
	return m.jobName
}

func (m *request) Data() map[string]interface{} {
	return m.data
}
