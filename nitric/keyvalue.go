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
	"fmt"

	"github.com/nitrictech/go-sdk/nitric/keyvalue"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type KvStorePermission string

const (
	KvStoreSet    KvStorePermission = "set"
	KvStoreGet    KvStorePermission = "get"
	KvStoreDelete KvStorePermission = "delete"
)

var KvStoreEverything []KvStorePermission = []KvStorePermission{KvStoreSet, KvStoreGet, KvStoreDelete}

type KvStore interface {
	// Allow requests the given permissions to the key/value store.
	Allow(KvStorePermission, ...KvStorePermission) (keyvalue.KvStoreClientIface, error)
}

type kvstore struct {
	name         string
	manager      *manager
	registerChan <-chan RegisterResult
}

// NewKv - Create a new Key/Value store resource
func NewKv(name string) *kvstore {
	kvstore := &kvstore{
		name:         name,
		manager:      defaultManager,
		registerChan: make(chan RegisterResult),
	}

	kvstore.registerChan = defaultManager.registerResource(&v1.ResourceDeclareRequest{
		Id: &v1.ResourceIdentifier{
			Type: v1.ResourceType_KeyValueStore,
			Name: name,
		},
		Config: &v1.ResourceDeclareRequest_KeyValueStore{
			KeyValueStore: &v1.KeyValueStoreResource{},
		},
	})

	return kvstore
}

func (k *kvstore) Allow(permission KvStorePermission, permissions ...KvStorePermission) (keyvalue.KvStoreClientIface, error) {
	allPerms := append([]KvStorePermission{permission}, permissions...)

	actions := []v1.Action{}
	for _, perm := range allPerms {
		switch perm {
		case KvStoreGet:
			actions = append(actions, v1.Action_KeyValueStoreRead)
		case KvStoreSet:
			actions = append(actions, v1.Action_KeyValueStoreWrite)
		case KvStoreDelete:
			actions = append(actions, v1.Action_KeyValueStoreDelete)
		default:
			return nil, fmt.Errorf("KvStorePermission %s unknown", perm)
		}
	}

	registerResult := <-k.registerChan

	if registerResult.Err != nil {
		return nil, registerResult.Err
	}

	err := k.manager.registerPolicy(registerResult.Identifier, actions...)
	if err != nil {
		return nil, err
	}

	return keyvalue.NewKvStoreClient(k.name)
}
