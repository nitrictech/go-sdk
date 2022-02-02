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

	nitricv1 "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/documents"
)

type CollectionPermission string

const (
	CollectionReading  CollectionPermission = "reading"
	CollectionWriting  CollectionPermission = "writing"
	CollectionDeleting CollectionPermission = "deleting"
)

var (
	CollectionEverything []CollectionPermission = []CollectionPermission{CollectionReading, CollectionWriting, CollectionDeleting}
)

func NewCollection(name string, permissions ...CollectionPermission) (documents.CollectionRef, error) {
	return run.NewCollection(name, permissions...)
}

func (m *manager) NewCollection(name string, permissions ...CollectionPermission) (documents.CollectionRef, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	colRes := &nitricv1.Resource{
		Type: nitricv1.ResourceType_Collection,
		Name: name,
	}

	dr := &nitricv1.ResourceDeclareRequest{
		Resource: colRes,
		Config: &nitricv1.ResourceDeclareRequest_Collection{
			Collection: &nitricv1.CollectionResource{},
		}}
	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []nitricv1.Action{}
	for _, perm := range permissions {
		switch perm {
		case CollectionReading:
			actions = append(actions, nitricv1.Action_CollectionDocumentRead, nitricv1.Action_CollectionList, nitricv1.Action_CollectionQuery)
		case CollectionWriting:
			actions = append(actions, nitricv1.Action_CollectionDocumentWrite)
		case CollectionDeleting:
			actions = append(actions, nitricv1.Action_CollectionDocumentDelete)
		default:
			return nil, fmt.Errorf("collectionPermission %s unknown", perm)
		}
	}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(colRes, actions))
	if err != nil {
		return nil, err
	}

	if m.docs == nil {
		m.docs, err = documents.New()
		if err != nil {
			return nil, err
		}
	}

	return m.docs.Collection(name), nil
}
