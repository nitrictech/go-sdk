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

package queues_examples

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/go-sdk/nitric/v1"
	"google.golang.org/grpc"
)

func TestQueues(t *testing.T) {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	// Create a mock storage service server...
	ctrl := gomock.NewController(t)
	ms := mock_v1.NewMockQueueServiceServer(ctrl)
	// Assert was called with the proper payload
	ms.EXPECT().Receive(gomock.Any(), gomock.Any()).Return(&v1.QueueReceiveResponse{
		Tasks: []*v1.NitricTask{
			{
				Id:      "1234",
				LeaseId: "1234",
			},
		},
	}, nil).Times(1)
	ms.EXPECT().SendBatch(gomock.Any(), gomock.Any()).Return(&v1.QueueSendBatchResponse{}, nil).Times(2)
	ms.EXPECT().Complete(gomock.Any(), gomock.Any()).Return(&v1.QueueCompleteResponse{}, nil).Times(1)
	// Start the gRPC server with the mock instance and await for it
	// to be called
	lis, _ := net.Listen("tcp", ":50051")

	v1.RegisterQueueServiceServer(grpcServer, ms)
	go grpcServer.Serve(lis)
	// call the functions to test
	send()
	receive()
	failed()
	// Cleanup
	grpcServer.Stop()
	lis.Close()
	ctrl.Finish()
}
