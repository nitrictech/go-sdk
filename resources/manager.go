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
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	multierror "github.com/missionMeteora/toolkit/errors"
	"google.golang.org/grpc"

	nitricv1 "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/documents"
	apierrors "github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/events"
	"github.com/nitrictech/go-sdk/api/queues"
	"github.com/nitrictech/go-sdk/api/secrets"
	"github.com/nitrictech/go-sdk/api/storage"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/faas"
)

type Starter interface {
	Start() error
}

type Manager interface {
	Run() error
	NewApi(name string, opts ...ApiOption) (Api, error)
	NewBucket(name string, permissions ...BucketPermission) (storage.Bucket, error)
	NewCollection(name string, permissions ...CollectionPermission) (documents.CollectionRef, error)
	NewSecret(name string, permissions ...SecretPermission) (secrets.SecretRef, error)
	NewQueue(name string, permissions ...QueuePermission) (queues.Queue, error)
	NewSchedule(name, rate string, handlers ...faas.EventMiddleware) error
	NewTopic(name string, permissions ...TopicPermission) (Topic, error)
}

type manager struct {
	blockers  map[string]Starter
	conn      grpc.ClientConnInterface
	connMutex sync.Mutex

	rsc      nitricv1.ResourceServiceClient
	evts     events.Events
	storage  storage.Storage
	docs     documents.Documents
	secrets  secrets.Secrets
	queues   queues.Queues
	builders map[string]faas.HandlerBuilder
}

var (
	run = &manager{
		blockers: map[string]Starter{},
		builders: map[string]faas.HandlerBuilder{},
	}
)

func New() Manager {
	return &manager{
		blockers: map[string]Starter{},
		builders: map[string]faas.HandlerBuilder{},
	}
}

func (m *manager) resourceServiceClient() (nitricv1.ResourceServiceClient, error) {
	m.connMutex.Lock()
	defer m.connMutex.Unlock()

	if m.conn == nil {
		conn, err := grpc.Dial(constants.NitricAddress(), constants.DefaultOptions()...)
		if err != nil {
			return nil, err
		}
		m.conn = conn
	}
	if m.rsc == nil {
		m.rsc = nitricv1.NewResourceServiceClient(m.conn)
	}
	return m.rsc, nil
}

func (m *manager) addStarter(name string, s Starter) {
	m.blockers[name] = s
}

func Run() error {
	return run.Run()
}

func (m *manager) Run() error {
	wg := sync.WaitGroup{}
	errList := &multierror.ErrorList{}

	for _, blocker := range m.blockers {
		wg.Add(1)
		go func(s Starter) {
			defer wg.Done()

			if err := s.Start(); err != nil {
				if IsBuildEnvirnonment() && isEOF(err) {
					// ignore the EOF error when running code-as-config.
					return
				}
				errList.Push(err)
			}
		}(blocker)
	}

	wg.Wait()

	return errList.Err()
}

func IsBuildEnvirnonment() bool {
	return strings.ToLower(os.Getenv("NITRIC_ENVIRONMENT")) == "build"
}

func isEOF(err error) bool {
	if err == nil {
		return false
	}
	var apiErr *apierrors.ApiError
	if errors.As(err, &apiErr) {
		err = apiErr.Unwrap()
	}
	if err == nil {
		return false
	}
	return errors.Is(err, io.EOF) || err.Error() == io.EOF.Error()
}
