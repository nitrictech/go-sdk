package storageclient

import (
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "go.nitric.io/go-sdk/interfaces/nitric/v1"
	mock_v1 "go.nitric.io/go-sdk/mocks"
)

var _ = Describe("Storageclient", func() {
	ctrl := gomock.NewController(GinkgoT())

	When("Put", func() {
		When("The bucket exists", func() {
			When("Putting a new item", func() {
				It("Should push the event", func() {
					mockStorageClient := mock_v1.NewMockStorageClient(ctrl)

					By("Calling Put with the expected inputs")
					mockStorageClient.EXPECT().Put(gomock.Any(), &v1.PutRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
						Body:       []byte{},
					} ).Return(&v1.PutReply{
						Success: true,
					}, nil)

					client := NewWithClient(mockStorageClient)
					err := client.Put("test-bucket", "test-key", []byte{})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})

				When("Success is false", func() {
					It("Should return an error", func() {
						mockStorageClient := mock_v1.NewMockStorageClient(ctrl)

						By("Calling Put with the expected inputs")
						mockStorageClient.EXPECT().Put(gomock.Any(), gomock.Any()).Return(&v1.PutReply{
							Success: false,
						}, nil)

						client := NewWithClient(mockStorageClient)
						err := client.Put("test-bucket", "test-key", []byte{})

						By("Returning an error")
						Expect(err).Should(HaveOccurred())
					})
				})
			})
		})

		When("An error is returned from the gRPC client", func() {
			It("Should return the error", func() {
				mockStorageClient := mock_v1.NewMockStorageClient(ctrl)

				By("Calling Put with the expected inputs")
				mockStorageClient.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil,
					fmt.Errorf("mock error"))

				client := NewWithClient(mockStorageClient)
				err := client.Put("test-bucket", "test-key", []byte{})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("Get", func() {
		When("The bucket exists", func() {
			When("The key exists", func() {
				It("Should retrieve the item", func() {
					mockStorageClient := mock_v1.NewMockStorageClient(ctrl)

					By("Calling Get with the expected inputs")
					mockStorageClient.EXPECT().Get(gomock.Any(), &v1.GetRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
						// FIXME: Using 'Reply' for storage but 'Request' for other services.
					}).Return(&v1.GetReply{
						Body: []byte{},
					}, nil)

					client := NewWithClient(mockStorageClient)
					item, err := client.Get("test-bucket", "test-key")

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

				By("Calling Get, which returns an error")
				mockStorageClient.EXPECT().Get(gomock.Any(), &v1.GetRequest{
					BucketName: "test-bucket",
					Key:        "test-key",
				}).Return(nil, fmt.Errorf("mock error"))

				client := NewWithClient(mockStorageClient)
				_, err := client.Get("test-bucket", "test-key")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})