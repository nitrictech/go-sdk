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

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
)

// Batch
type Batch interface {
	// Job - Retrieve a Job reference
	Job(name string) Job
}

type batchImpl struct {
	batchClient v1.BatchClient
}

func (s *batchImpl) Job(name string) Job {
	// Just return the straight job reference
	// we can fail if the job does not exist
	return &jobImpl{
		name:        name,
		batchClient: s.batchClient,
	}
}

// New - Construct a new Batch Client with default options
func New() (Batch, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.NitricDialTimeout())
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		return nil, errors.NewWithCause(codes.Unavailable, "Unable to dial Batch service", err)
	}

	tc := v1.NewBatchClient(conn)

	return &batchImpl{
		batchClient: tc,
	}, nil
}
