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
	"errors"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	"github.com/nitrictech/protoutils"
)

var _ = Describe("File", func() {
	var (
		ctrl            *gomock.Controller
		mockBatchClient *mock_v1.MockBatchClient
		j               Job
		jobName         string
		ctx             context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockBatchClient = mock_v1.NewMockBatchClient(ctrl)

		jobName = "test-job"
		j = &jobImpl{
			name:        jobName,
			batchClient: mockBatchClient,
		}

		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Name()", func() {
		It("should have the same job name as the one provided", func() {
			_jobName := j.Name()
			Expect(_jobName).To(Equal(jobName))
		})
	})

	Describe("Submit()", func() {
		var dataToBeSubmitted map[string]interface{}

		BeforeEach(func() {
			dataToBeSubmitted = map[string]interface{}{
				"data": "hello world",
			}
		})

		When("the gRPC Read operation is successful", func() {
			BeforeEach(func() {
				payloadStruct, err := protoutils.NewStruct(dataToBeSubmitted)
				Expect(err).ToNot(HaveOccurred())

				mockBatchClient.EXPECT().SubmitJob(gomock.Any(), &v1.SubmitJobRequest{
					Name: jobName,
					Data: &v1.JobData{
						Data: &v1.JobData_Struct{
							Struct: payloadStruct,
						},
					},
				}).Return(
					&v1.SubmitJobResponse{},
					nil).Times(1)
			})

			It("should not return error", func() {
				err := j.Submit(ctx, dataToBeSubmitted)

				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("the grpc server returns an error", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "Internal Error"

				By("the gRPC server returning an error")
				mockBatchClient.EXPECT().SubmitJob(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New(errorMsg),
				).Times(1)
			})

			It("should return the passed error", func() {
				err := j.Submit(ctx, dataToBeSubmitted)

				By("returning error with expected message")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
			})
		})
	})
})
