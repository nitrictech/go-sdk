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

package faas

import (
	"fmt"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

// APIS

type ApiWorkerOptions struct {
	ApiName          string
	Path             string
	Security         map[string][]string
	SecurityDisabled bool
}

// SCHEDULES

type Frequency string //= "days" | "hours" | "minutes";

var Frequencies = []Frequency{"days", "hours", "minutes"}

type RateWorkerOptions struct {
	Description string
	Rate        int
	Frequency   Frequency
}

type CronWorkerOptions struct {
	Description string
	Cron        string
}

// TOPICS

type SubscriptionWorkerOptions struct {
	Topic string
}

// BUCKET NOTIFICATIONS

type NotificationType string

var NotificationTypes = []NotificationType{WriteNotification, DeleteNotification}

const (
	WriteNotification  NotificationType = "write"
	DeleteNotification NotificationType = "delete"
)

type BucketNotificationWorkerOptions struct {
	Bucket                   string
	NotificationType         NotificationType
	NotificationPrefixFilter string
}

func (b *BucketNotificationWorkerOptions) notificationTypeToWire() (v1.BucketNotificationType, error) {
	switch b.NotificationType {
	case WriteNotification:
		return v1.BucketNotificationType_Created, nil
	case DeleteNotification:
		return v1.BucketNotificationType_Deleted, nil
	default:
		return -1, fmt.Errorf("notification type %s is unsupported", b.NotificationType)
	}
}

// WEBSOCKETS

type WebsocketEventType string

var WebsocketEventTypes = []WebsocketEventType{WebsocketConnect, WebsocketDisconnect, WebsocketMessage}

const (
	WebsocketConnect    WebsocketEventType = "connect"
	WebsocketDisconnect WebsocketEventType = "disconnect"
	WebsocketMessage    WebsocketEventType = "message"
)

type WebsocketWorkerOptions struct {
	Socket    string
	EventType WebsocketEventType
}

func (w *WebsocketWorkerOptions) eventTypeToWire() (v1.WebsocketEvent, error) {
	switch w.EventType {
	case WebsocketConnect:
		return v1.WebsocketEvent_Connect, nil
	case WebsocketDisconnect:
		return v1.WebsocketEvent_Disconnect, nil
	case WebsocketMessage:
		return v1.WebsocketEvent_Message, nil
	default:
		return -1, fmt.Errorf("websocket type %s is unsupported", w.EventType)
	}
}
