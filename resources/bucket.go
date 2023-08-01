// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"context"
	"fmt"

	"github.com/nitrictech/go-sdk/api/storage"
	"github.com/nitrictech/go-sdk/faas"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

type BucketPermission string

type bucket struct {
	name string
	manager Manager
}

type Bucket interface {
	With(permissions ...BucketPermission) (storage.Bucket, error)
	On(faas.NotificationType, string, ...faas.BucketNotificationMiddleware)
}

const (
	BucketReading  BucketPermission = "reading"
	BucketWriting  BucketPermission = "writing"
	BucketDeleting BucketPermission = "deleting"
)

var BucketEverything []BucketPermission = []BucketPermission{BucketReading, BucketWriting, BucketDeleting}

// NewBucket register this bucket as a required resource for the calling function/container and
// register the permissions required by the currently scoped function for this resource.
func NewBucket(name string) Bucket {
	return &bucket{
		name: name,
		manager: defaultManager,
	}
}

func (b *bucket) With(permissions ...BucketPermission) (storage.Bucket, error) {
	return defaultManager.newBucket(b.name, permissions...)
}

func (m *manager) newBucket(name string, permissions ...BucketPermission) (storage.Bucket, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	res := &v1.Resource{
		Type: v1.ResourceType_Bucket,
		Name: name,
	}

	dr := &v1.ResourceDeclareRequest{
		Resource: res,
		Config: &v1.ResourceDeclareRequest_Bucket{
			Bucket: &v1.BucketResource{},
		},
	}
	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []v1.Action{}
	for _, perm := range permissions {
		switch perm {
		case BucketReading:
			actions = append(actions, v1.Action_BucketFileGet, v1.Action_BucketFileList)
		case BucketWriting:
			actions = append(actions, v1.Action_BucketFilePut)
		case BucketDeleting:
			actions = append(actions, v1.Action_BucketFileDelete)
		default:
			return nil, fmt.Errorf("bucketPermission %s unknown", perm)
		}
	}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(res, actions))
	if err != nil {
		return nil, err
	}

	if m.storage == nil {
		m.storage, err = storage.New()
		if err != nil {
			return nil, err
		}
	}

	return m.storage.Bucket(name), nil
}

func (b *bucket) On(notificationType faas.NotificationType, notificationPrefixFilter string, middleware ...faas.BucketNotificationMiddleware) {	
	f := faas.New()

	f.BucketNotification(middleware...)
	f.WithBucketNotificationWorkerOptions(faas.BucketNotificationWorkerOptions{Bucket: b.name, NotificationType: notificationType, NotificationPrefixFilter: notificationPrefixFilter})

	b.manager.addWorker(fmt.Sprintf("bucket:notification %s %s/%s", b.name, notificationType, notificationPrefixFilter), f)
}