// Copyright 2023 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// XXX: File requests currently not implemented
// type FileRequest interface {
// 	Bucket() *Bucket
// 	NotificationType() EventType
// }

// type fileRequestImpl struct {
// 	bucket           Bucket
// 	notificationType EventType
// }

// func (f *fileRequestImpl) Bucket() Bucket {
// 	return f.bucket
// }

// func (f *fileRequestImpl) NotificationType() EventType {
// 	return f.notificationType
// }
