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

package nitric

import (
	"context"
	"fmt"
	"strings"

	"github.com/nitrictech/go-sdk/api/storage"
	"github.com/nitrictech/go-sdk/handler"
	"github.com/nitrictech/go-sdk/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

type BucketPermission string

type bucket struct {
	name    string
	manager Manager
}

type Bucket interface {
	Allow(BucketPermission, ...BucketPermission) (storage.Bucket, error)
	On(handler.BlobEventType, string, ...handler.BlobEventMiddleware)
}

const (
	BucketRead   BucketPermission = "read"
	BucketWrite  BucketPermission = "write"
	BucketDelete BucketPermission = "delete"
)

var BucketEverything []BucketPermission = []BucketPermission{BucketRead, BucketWrite, BucketDelete}

// NewBucket register this bucket as a required resource for the calling function/container and
// register the permissions required by the currently scoped function for this resource.
func NewBucket(name string) Bucket {
	return &bucket{
		name:    name,
		manager: defaultManager,
	}
}

func (b *bucket) Allow(permission BucketPermission, permissions ...BucketPermission) (storage.Bucket, error) {
	allPerms := append([]BucketPermission{permission}, permissions...)

	return defaultManager.newBucket(b.name, allPerms...)
}

func (m *manager) newBucket(name string, permissions ...BucketPermission) (storage.Bucket, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	res := &v1.ResourceIdentifier{
		Type: v1.ResourceType_Bucket,
		Name: name,
	}

	dr := &v1.ResourceDeclareRequest{
		Id: res,
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
		case BucketRead:
			actions = append(actions, v1.Action_BucketFileGet, v1.Action_BucketFileList)
		case BucketWrite:
			actions = append(actions, v1.Action_BucketFilePut)
		case BucketDelete:
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

func (b *bucket) On(notificationType handler.BlobEventType, notificationPrefixFilter string, middleware ...handler.BlobEventMiddleware) {
	var blobEventType storagepb.BlobEventType;
	switch notificationType {
	case handler.WriteNotification:
		blobEventType = storagepb.BlobEventType_Created
	case handler.DeleteNotification:
		blobEventType = storagepb.BlobEventType_Deleted
	}

	registrationRequest := &storagepb.RegistrationRequest{
		BucketName: b.name,
		BlobEventType: blobEventType,
		KeyPrefixFilter: notificationPrefixFilter,
	}
	
	composedHandler := handler.ComposeBlobEventMiddleware(middleware...)

	opts := &workers.BlobEventWorkerOpts{
		RegistrationRequest: registrationRequest,
		Middleware: composedHandler,
	}

	worker := workers.NewBlobEventWorker(opts)
	
	b.manager.addWorker("bucketNotification:"+ strings.Join([]string{
		b.name, notificationPrefixFilter, string(notificationType),
	},"-"), worker)
}
