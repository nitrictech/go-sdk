package storageclient

import (
	"fmt"

	"github.com/golang/mock/gomock"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storageclient", func() {
	ctrl := gomock.NewController(GinkgoT())

	When("Write", func() {
		When("The bucket exists", func() {
			When("Writeting a new item", func() {
				It("Should push the event", func() {
					mockStorageClient := mock_v1.NewMockStorageClient(ctrl)

					By("Calling Write with the expected inWrites")
					mockStorageClient.EXPECT().Write(gomock.Any(), &v1.StorageWriteRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
						Body:       []byte{},
					}).Return(&v1.StorageWriteResponse{}, nil)

					client := NewWithClient(mockStorageClient)
					err := client.Write("test-bucket", "test-key", []byte{})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		When("An error is returned from the gRPC client", func() {
			It("Should return the error", func() {
				mockStorageClient := mock_v1.NewMockStorageClient(ctrl)

				By("Calling Write with the expected inWrites")
				mockStorageClient.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil,
					fmt.Errorf("mock error"))

				client := NewWithClient(mockStorageClient)
				err := client.Write("test-bucket", "test-key", []byte{})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("Read", func() {
		When("The bucket exists", func() {
			When("The key exists", func() {
				It("Should retrieve the item", func() {
					mockStorageClient := mock_v1.NewMockStorageClient(ctrl)

					By("Calling Read with the expected inWrites")
					mockStorageClient.EXPECT().Read(gomock.Any(), &v1.StorageReadRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
						// FIXME: Using 'Reply' for storage but 'Request' for other services.
					}).Return(&v1.StorageReadResponse{
						Body: []byte{},
					}, nil)

					client := NewWithClient(mockStorageClient)
					item, err := client.Read("test-bucket", "test-key")

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the item")
					Expect(item).To(Equal([]byte{}))
				})
			})

			When("The key doesn't exists", func() {
				// TODO: handle not found gRPC error
			})
		})

		When("An error is returned from the gRPC client", func() {
			It("Should return an error", func() {
				mockStorageClient := mock_v1.NewMockStorageClient(ctrl)

				By("Calling Read, which returns an error")
				mockStorageClient.EXPECT().Read(gomock.Any(), &v1.StorageReadRequest{
					BucketName: "test-bucket",
					Key:        "test-key",
				}).Return(nil, fmt.Errorf("mock error"))

				client := NewWithClient(mockStorageClient)
				_, err := client.Read("test-bucket", "test-key")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
