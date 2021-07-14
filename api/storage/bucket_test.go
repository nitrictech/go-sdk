package storage

import (
	"github.com/golang/mock/gomock"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bucket", func() {
	ctrl := gomock.NewController(GinkgoT())

	Context("File", func() {
		When("creating a new File reference", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)

			bucket := &bucketImpl{
				name: "test-bucket",
				sc:   mockStorage,
			}

			object := bucket.File("test-object")
			objectI, ok := object.(*fileImpl)

			It("should return an objectImpl instance", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provided file name", func() {
				Expect(objectI.key).To(Equal("test-object"))
			})

			It("should share the buckets name", func() {
				Expect(objectI.bucket).To(Equal("test-bucket"))
			})

			It("should share the bucket references gRPC client", func() {
				Expect(objectI.sc).To(Equal(mockStorage))
			})
		})
	})
})
