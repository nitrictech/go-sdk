package topics

type Request interface {
	TopicName() string
	Message() map[string]interface{}
}

type requestImpl struct {
	topicName string
	message   map[string]interface{}
}

func (m *requestImpl) TopicName() string {
	return m.topicName
}

func (m *requestImpl) Message() map[string]interface{} {
	return m.message
}
