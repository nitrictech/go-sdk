package authclient

import (
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("Authclient", func() {
	ctrl := gomock.NewController(GinkgoT())

	When("CreateUser", func() {
		When("The user is created by the gRPC server", func() {
			It("Should call gRPC CreateUser", func() {
				mockGRPCAuthClient := mock_v1.NewMockAuthClient(ctrl)

				By("Calling CreateUser with the expected inputs")
				mockGRPCAuthClient.EXPECT().
					CreateUser(
						gomock.Any(),
						&v1.CreateUserRequest{
							Tenant:   "test-tenant",
							Id:       "testid",
							Email:    "test@example.com",
							Password: "testpassword",
						},
					)

				client := NewWithClient(mockGRPCAuthClient)
				err := client.CreateUser("test-tenant", "testid", "test@example.com", "testpassword")

				By("No returning an error")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})


		When("An error is returned from the gRPC server", func () {
			It("Should call gRPC CreateUser", func() {
				mockGRPCAuthClient := mock_v1.NewMockAuthClient(ctrl)

				By("Calling CreateUser")
				mockGRPCAuthClient.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("mock error"))

				client := NewWithClient(mockGRPCAuthClient)
				err := client.CreateUser("test-tenant", "testid", "test@example.com", "testpassword")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
