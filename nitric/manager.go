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

	"github.com/nitrictech/go-sdk/api/batch"
	apierrors "github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/keyvalue"
	"github.com/nitrictech/go-sdk/api/queues"
	"github.com/nitrictech/go-sdk/api/secrets"
	"github.com/nitrictech/go-sdk/api/storage"
	"github.com/nitrictech/go-sdk/api/topics"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

// Manager is the top level object that resources are created on.
type Manager interface {
	Run() error
	addWorker(name string, s workers.Worker)
	resourceServiceClient() (v1.ResourcesClient, error)
	registerResource(request *v1.ResourceDeclareRequest) <-chan RegisterResult
	registerPolicy(res *v1.ResourceIdentifier, actions ...v1.Action) (*manager, error)
}

type RegisterResult struct {
	Identifier *v1.ResourceIdentifier
	Err        error
}

type manager struct {
	workers   map[string]workers.Worker
	conn      grpc.ClientConnInterface
	connMutex sync.Mutex

	rsc      v1.ResourcesClient
	topics   topics.Topics
	batch    batch.Batch
	storage  storage.Storage
	secrets  secrets.Secrets
	queues   queues.Queues
	kvstores keyvalue.KeyValue
}

var defaultManager = New()

// New is used to create the top level resource manager.
// Note: this is not required if you are using
// resources.NewApi() and the like. These use a default manager instance.
func New() Manager {
	return &manager{
		workers: map[string]workers.Worker{},
	}
}

func (m *manager) addWorker(name string, s workers.Worker) {
	m.workers[name] = s
}

func (m *manager) resourceServiceClient() (v1.ResourcesClient, error) {
	m.connMutex.Lock()
	defer m.connMutex.Unlock()

	if m.conn == nil {
		ctx, cancel := context.WithTimeout(context.Background(), constants.NitricDialTimeout())
		defer cancel()

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
		m.rsc = v1.NewResourcesClient(m.conn)
	}
	return m.rsc, nil
}

func (m *manager) registerResource(request *v1.ResourceDeclareRequest) <-chan RegisterResult {
	registerResourceChan := make(chan RegisterResult)

	go func() {
		rsc, err := m.resourceServiceClient()
		if err != nil {
			registerResourceChan <- RegisterResult{
				Err:        err,
				Identifier: nil,
			}

			return
		}

		_, err = rsc.Declare(context.Background(), request)
		if err != nil {
			registerResourceChan <- RegisterResult{
				Err:        err,
				Identifier: nil,
			}

			return
		}

		registerResourceChan <- RegisterResult{
			Err:        nil,
			Identifier: request.Id,
		}
	}()

	return registerResourceChan
}

func (m *manager) registerPolicy(res *v1.ResourceIdentifier, actions ...v1.Action) (*manager, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return m, err
	}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(res, actions))
	if err != nil {
		return m, err
	}

	return m, nil
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
		go func(s workers.Worker) {
			defer wg.Done()

			if err := s.Start(context.TODO()); err != nil {
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
