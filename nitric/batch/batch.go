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

package batch

import (
	"fmt"

	"github.com/nitrictech/go-sdk/internal/handlers"
	"github.com/nitrictech/go-sdk/nitric/workers"
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

// JobPermission defines the available permissions on a job
type JobPermission string

type Handler = handlers.Handler[Ctx]

const (
	// JobSubmit is required to call Submit on a job.
	JobSubmit JobPermission = "submit"
)

type JobReference interface {
	// Allow requests the given permissions to the job.
	Allow(permission JobPermission, permissions ...JobPermission) *BatchClient

	// Handler will register and start the job task handler that will be called for all task submitted to this job.
	// Valid function signatures for middleware are:
	//
	//	func()
	//	func() error
	//	func(*batch.Ctx)
	//	func(*batch.Ctx) error
	//	Handler[batch.Ctx]
	Handler(handler interface{}, options ...HandlerOption)
}

type jobReference struct {
	name         string
	manager      *workers.Manager
	registerChan <-chan workers.RegisterResult
}

// NewJob creates a new job resource with the give name.
func NewJob(name string) JobReference {
	job := &jobReference{
		name:    name,
		manager: workers.GetDefaultManager(),
	}

	job.registerChan = job.manager.RegisterResource(&v1.ResourceDeclareRequest{
		Id: &v1.ResourceIdentifier{
			Type: v1.ResourceType_Job,
			Name: name,
		},
		Config: &v1.ResourceDeclareRequest_Job{
			Job: &v1.JobResource{},
		},
	})

	return job
}

func (j *jobReference) Allow(permission JobPermission, permissions ...JobPermission) *BatchClient {
	allPerms := append([]JobPermission{permission}, permissions...)

	actions := []v1.Action{}
	for _, perm := range allPerms {
		switch perm {
		case JobSubmit:
			actions = append(actions, v1.Action_JobSubmit)
		default:
			panic(fmt.Errorf("JobPermission %s unknown", perm))
		}
	}

	registerResult := <-j.registerChan
	if registerResult.Err != nil {
		panic(registerResult.Err)
	}

	err := j.manager.RegisterPolicy(registerResult.Identifier, actions...)
	if err != nil {
		panic(err)
	}

	client, err := NewBatchClient(j.name)
	if err != nil {
		panic(err)
	}

	return client
}

func (j *jobReference) Handler(handler interface{}, opts ...HandlerOption) {
	options := &handlerOptions{}

	for _, opt := range opts {
		opt(options)
	}

	registrationRequest := &batchpb.RegistrationRequest{
		JobName:      j.name,
		Requirements: &batchpb.JobResourceRequirements{},
	}

	if options.cpus != nil {
		registrationRequest.Requirements.Cpus = *options.cpus
	}

	if options.memory != nil {
		registrationRequest.Requirements.Memory = *options.memory
	}

	if options.gpus != nil {
		registrationRequest.Requirements.Gpus = *options.gpus
	}

	typedHandler, err := handlers.HandlerFromInterface[Ctx](handler)
	if err != nil {
		panic(err)
	}

	jobOpts := &jobWorkerOpts{
		RegistrationRequest: registrationRequest,
		Handler:             typedHandler,
	}

	worker := newJobWorker(jobOpts)
	j.manager.AddWorker("JobWorker:"+j.name, worker)
}
