package events

type Event struct {
	Payload     map[string]interface{}
	PayloadType string
	ID          string
}
