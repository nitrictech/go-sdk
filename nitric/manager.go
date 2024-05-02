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
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	multierror "github.com/missionMeteora/toolkit/errors"
	"google.golang.org/grpc"

	apierrors "github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/events"
	"github.com/nitrictech/go-sdk/api/queues"
	"github.com/nitrictech/go-sdk/api/secrets"
	"github.com/nitrictech/go-sdk/api/storage"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/faas"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type Starter interface {
	Start() error
}

// Manager is the top level object that resources are created on.
type Manager interface {
	Run() error
	addWorker(name string, s Starter)
	addBuilder(name string, builder faas.HandlerBuilder)
	getBuilder(name string) faas.HandlerBuilder
	resourceServiceClient() (v1.ResourceServiceClient, error)

	newApi(name string, opts ...ApiOption) (Api, error)
	newBucket(name string, permissions ...BucketPermission) (storage.Bucket, error)
	newSecret(name string, permissions ...SecretPermission) (secrets.SecretRef, error)
	newQueue(name string, permissions ...QueuePermission) (queues.Queue, error)
	newSchedule(name string) Schedule
	newTopic(name string, permissions ...TopicPermission) (Topic, error)
	newWebsocket(socket string) (Websocket, error)
}

type manager struct {
	workers   map[string]Starter
	conn      grpc.ClientConnInterface
	connMutex sync.Mutex

	rsc      v1.ResourceServiceClient
	evts     events.Events
	storage  storage.Storage
	secrets  secrets.Secrets
	queues   queues.Queues
	builders map[string]faas.HandlerBuilder
}

var (
	defaultManager = New()
	traceInit      = sync.Once{}
)

// New is used to create the top level resource manager.
// Note: this is not required if you are using
// resources.NewApi() and the like. These use a default manager instance.
func New() Manager {
	return &manager{
		workers:  map[string]Starter{},
		builders: map[string]faas.HandlerBuilder{},
	}
}

// Gets an existing builder or returns a new handler builder
func (m *manager) getBuilder(name string) faas.HandlerBuilder {
	return m.builders[name]
}

func (m *manager) addBuilder(name string, builder faas.HandlerBuilder) {
	m.builders[name] = builder
}

func (m *manager) addWorker(name string, s Starter) {
	m.workers[name] = s
}

func (m *manager) resourceServiceClient() (v1.ResourceServiceClient, error) {
	m.connMutex.Lock()
	defer m.connMutex.Unlock()

	if m.conn == nil {
		ctx, _ := context.WithTimeout(context.TODO(), constants.NitricDialTimeout())

		conn, err := grpc.DialContext(
			ctx,
			constants.NitricAddress(),
			constants.DefaultOptions()...,
		)
		if err != nil {
			return nil, err
		}
		m.conn = conn
	}
	if m.rsc == nil {
		m.rsc = v1.NewResourceServiceClient(m.conn)
	}
	return m.rsc, nil
}

// Run will run the function and callback the required handlers when these events are received.
func Run() error {
	return defaultManager.Run()
}

func (m *manager) Run() error {
	wg := sync.WaitGroup{}
	errList := &multierror.ErrorList{}

	for _, worker := range m.workers {
		wg.Add(1)
		go func(s Starter) {
			defer wg.Done()

			if err := s.Start(); err != nil {
				if isBuildEnvirnonment() && isEOF(err) {
					// ignore the EOF error when running code-as-config.
					return
				}

				errList.Push(err)
			}
		}(worker)
	}

	wg.Wait()

	return errList.Err()
}

// IsBuildEnvirnonment will return true if the code is running during config discovery.
func isBuildEnvirnonment() bool {
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
