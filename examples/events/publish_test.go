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

package events_examples

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "github.com/nitrictech/apis/go/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	"google.golang.org/grpc"
)

func TestPublishTopic(t *testing.T) {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	// Create a mock storage service server...
	ctrl := gomock.NewController(t)
	mec := mock_v1.NewMockEventServiceServer(ctrl)
	mtc := mock_v1.NewMockTopicServiceServer(ctrl)
	// Assert was called with the proper payload
	mec.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(&v1.EventPublishResponse{}, nil).Times(2)

	// Start the gRPC server with the mock instance and await for it
	// to be called
	lis, _ := net.Listen("tcp", ":50051")

	v1.RegisterEventServiceServer(grpcServer, mec)
	v1.RegisterTopicServiceServer(grpcServer, mtc)
	go grpcServer.Serve(lis)
	// call the function to test
	publishEvent()
	publishEventId()
	// Cleanup
	grpcServer.Stop()
	lis.Close()
	ctrl.Finish()
}
