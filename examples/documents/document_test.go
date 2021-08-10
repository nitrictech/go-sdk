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

package documents_examples

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	"google.golang.org/grpc"
)

func TestDeleteDocument(t *testing.T) {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	// Create a mock storage service server...
	ctrl := gomock.NewController(t)
	ms := mock_v1.NewMockDocumentServiceServer(ctrl)
	// Assert was called with the proper payload
	ms.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&v1.DocumentDeleteResponse{}, nil).Times(1)
	ms.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&v1.DocumentGetResponse{}, nil).Times(1)
	ms.EXPECT().Set(gomock.Any(), gomock.Any()).Return(&v1.DocumentSetResponse{}, nil).Times(1)
	ms.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&v1.DocumentQueryResponse{}, nil).Times(7)
	ms.EXPECT().QueryStream(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	// Start the gRPC server with the mock instance and await for it
	// to be called
	lis, _ := net.Listen("tcp", ":50051")

	v1.RegisterDocumentServiceServer(grpcServer, ms)
	go grpcServer.Serve(lis)
	// call the function to test
	delete()
	get()
	set()
	query()
	queryFilter()
	queryLimit()
	subDocQuery()
	subColQuery()
	pagedResults()
	stream()
	// Cleanup
	grpcServer.Stop()
	lis.Close()
	ctrl.Finish()
}
