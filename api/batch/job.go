package batch

import (
	"context"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	"github.com/nitrictech/protoutils"
)

type SubmitOption = func(*v1.JobSubmitRequest)

type Job interface {
	// Name returns the Job name.
	Name() string

	// Submit the provided job to the job service.
	Submit(context.Context, map[string]interface{}, ...SubmitOption) error
}

type jobImpl struct {
	name        string
	batchClient v1.BatchClient
}

func (s *jobImpl) Name() string {
	return s.name
}

func (s *jobImpl) Submit(ctx context.Context, data map[string]interface{}, opts ...SubmitOption) error {
	// Convert data to Protobuf Struct
	dataStruct, err := protoutils.NewStruct(data)
	if err != nil {
		return errors.NewWithCause(codes.InvalidArgument, "Topic.Publish", err)
	}

	req := &v1.JobSubmitRequest{
		JobName: s.name,
		Data: &v1.JobData{
			Data: &v1.JobData_Struct{
				Struct: dataStruct,
			},
		},
	}

	// Apply options to the request
	for _, opt := range opts {
		opt(req)
	}

	_, err = s.batchClient.SubmitJob(ctx, req)
	if err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}
