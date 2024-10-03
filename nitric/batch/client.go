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
	"context"

	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	"github.com/nitrictech/protoutils"
)

// Batch
type BatchClientIn interface {
	// Name returns the Job name.
	Name() string

	// Submit will submit the provided request to the job.
	Submit(ctx context.Context, data map[string]interface{}) error
}

type BatchClient struct {
	name        string
	batchClient v1.BatchClient
}

func (s *BatchClient) Name() string {
	return s.name
}

func (s *BatchClient) Submit(ctx context.Context, data map[string]interface{}) error {
	dataStruct, err := protoutils.NewStruct(data)
	if err != nil {
		return errors.NewWithCause(codes.InvalidArgument, "Batch.Submit", err)
	}

	// Create the request
	req := &v1.JobSubmitRequest{
		JobName: s.name,
		Data: &v1.JobData{
			Data: &v1.JobData_Struct{
				Struct: dataStruct,
			},
		},
	}

	// Submit the request
	_, err = s.batchClient.SubmitJob(ctx, req)
	if err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

func NewBatchClient(name string) (*BatchClient, error) {
	conn, err := grpc.NewClient(constants.NitricAddress(), constants.DefaultOptions()...)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"NewBatchClient: unable to reach nitric server",
			err,
		)
	}

	batchClient := v1.NewBatchClient(conn)

	return &BatchClient{
		name:        name,
		batchClient: batchClient,
	}, nil
}
