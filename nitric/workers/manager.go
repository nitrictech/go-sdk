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

package workers

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	multierror "github.com/missionMeteora/toolkit/errors"

	grpcx "github.com/nitrictech/go-sdk/internal/grpc"
	apierrors "github.com/nitrictech/go-sdk/nitric/errors"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type RegisterResult struct {
	Identifier *v1.ResourceIdentifier
	Err        error
}

type Manager struct {
	workers map[string]StreamWorker

	rsc v1.ResourcesClient
}

var defaultManager = New()

func GetDefaultManager() *Manager {
	return defaultManager
}

// New is used to create the top level resource manager.
// Note: this is not required if you are using
// resources.NewApi() and the like. These use a default manager instance.
func New() *Manager {
	return &Manager{
		workers: map[string]StreamWorker{},
	}
}

func (m *Manager) AddWorker(name string, s StreamWorker) {
	m.workers[name] = s
}

func (m *Manager) resourceServiceClient() (v1.ResourcesClient, error) {
	conn, err := grpcx.GetConnection()
	if err != nil {
		return nil, err
	}

	if m.rsc == nil {
		m.rsc = v1.NewResourcesClient(conn)
	}
	return m.rsc, nil
}

func (m *Manager) RegisterResource(request *v1.ResourceDeclareRequest) <-chan RegisterResult {
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

func functionResourceDeclareRequest(subject *v1.ResourceIdentifier, actions []v1.Action) *v1.ResourceDeclareRequest {
	return &v1.ResourceDeclareRequest{
		Id: &v1.ResourceIdentifier{
			Type: v1.ResourceType_Policy,
		},
		Config: &v1.ResourceDeclareRequest_Policy{
			Policy: &v1.PolicyResource{
				Principals: []*v1.ResourceIdentifier{
					{
						Type: v1.ResourceType_Service,
					},
				},
				Actions:   actions,
				Resources: []*v1.ResourceIdentifier{subject},
			},
		},
	}
}

func (m *Manager) RegisterPolicy(res *v1.ResourceIdentifier, actions ...v1.Action) error {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return err
	}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(res, actions))
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	errList := &multierror.ErrorList{}

	for _, worker := range m.workers {
		wg.Add(1)
		go func(s StreamWorker) {
			defer wg.Done()

			if err := s.Start(ctx); err != nil {
				if isBuildEnvironment() && isEOF(err) {
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

// IsBuildEnvironment will return true if the code is running during config discovery.
func isBuildEnvironment() bool {
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
