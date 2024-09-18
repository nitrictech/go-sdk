package storage

type EventType string

var EventTypes = []EventType{WriteNotification, DeleteNotification}

const (
	WriteNotification  EventType = "write"
	DeleteNotification EventType = "delete"
)

type Request interface {
	Key() string
	NotificationType() EventType
}

type requestImpl struct {
	key              string
	notificationType EventType
}

func (b *requestImpl) Key() string {
	return b.key
}

func (b *requestImpl) NotificationType() EventType {
	return b.notificationType
}

// File Event

type FileRequest interface {
	Bucket() *Bucket
	NotificationType() EventType
}

type fileRequestImpl struct {
	bucket           Bucket
	notificationType EventType
}

func (f *fileRequestImpl) Bucket() Bucket {
	return f.bucket
}

func (f *fileRequestImpl) NotificationType() EventType {
	return f.notificationType
}
