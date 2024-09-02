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

	"github.com/nitrictech/go-sdk/api/batch"
	"github.com/nitrictech/go-sdk/handler"
	"github.com/nitrictech/go-sdk/workers"
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

// JobPermission defines the available permissions on a job
type JobPermission string

const (
	// JobSubmit is required to call Submit on a job.
	JobSubmit JobPermission = "submit"
)

type Job interface {
	batch.Job
}

type JobResourceRequirements struct {
	Cpus   float32
	Memory int64
	Gpus   int64
}

// The resource declaration, not the runtime interact-able object
type JobDefinition interface {
	Allow(JobPermission, ...JobPermission) (batch.Job, error)

	// Handler - Set the job handler function
	Handler(JobResourceRequirements, ...handler.JobMiddleware)
}

type job struct {
	batch.Job

	manager Manager
}

type definableJob struct {
	name         string
	manager      Manager
	registerChan <-chan RegisterResult
}

func NewJob(name string) JobDefinition {
	job := &definableJob{
		name:    name,
		manager: defaultManager,
	}

	job.registerChan = defaultManager.registerResource(&v1.ResourceDeclareRequest{
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

func (j *definableJob) Allow(permission JobPermission, permissions ...JobPermission) (batch.Job, error) {
	allPerms := append([]JobPermission{permission}, permissions...)

	actions := []v1.Action{}
	for _, perm := range allPerms {
		switch perm {
		case JobSubmit:
			actions = append(actions, v1.Action_JobSubmit)
		default:
			return nil, fmt.Errorf("JobPermission %s unknown", perm)
		}
	}

	registerResult := <-j.registerChan
	if registerResult.Err != nil {
		return nil, registerResult.Err
	}

	m, err := j.manager.registerPolicy(registerResult.Identifier, actions...)
	if err != nil {
		return nil, err
	}

	if m.batch == nil {
		evts, err := batch.New()
		if err != nil {
			return nil, err
		}
		m.batch = evts
	}

	return &job{
		Job:     m.batch.Job(j.name),
		manager: m,
	}, nil
}

func (j *definableJob) Handler(reqs JobResourceRequirements, middleware ...handler.JobMiddleware) {
	registrationRequest := &batchpb.RegistrationRequest{
		JobName: j.name,
		Requirements: &batchpb.JobResourceRequirements{
			Cpus:   reqs.Cpus,
			Memory: reqs.Memory,
			Gpus:   reqs.Gpus,
		},
	}
	composeHandler := handler.ComposeJobMiddleware(middleware...)

	opts := &workers.JobWorkerOpts{
		RegistrationRequest: registrationRequest,
		Middleware:          composeHandler,
	}

	worker := workers.NewJobWorker(opts)
	j.manager.addWorker("JobWorker:"+j.name, worker)
}
