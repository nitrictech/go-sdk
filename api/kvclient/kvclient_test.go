package kvclient_test

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/golang/mock/gomock"
	"github.com/nitrictech/go-sdk/api/kvclient"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ = Describe("KVclient", func() {
	ctrl := gomock.NewController(GinkgoT())

	When("GetKey", func() {
		When("The collection exists", func() {
			When("The key exists", func() {
				It("Should retrieve the value", func() {
					mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

					mockValue, _ := structpb.NewStruct(map[string]interface{}{
						"test": 123,
					})

					By("Calling GetKey with the expected inputs")
					mockKvClient.EXPECT().
						Get(
							gomock.Any(),
							&v1.KeyValueGetRequest{
								Collection: "test-collection",
								Key:        "test-key",
							},
						).Return(&v1.KeyValueGetResponse{
						Value: mockValue,
					}, nil)

					client := kvclient.NewWithClient(mockKvClient)
					val, err := client.GetKey("test-collection", "test-key")

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the document")
					Expect(int(val["test"].(float64))).To(Equal(123))
				})
			})
			When("The key doesn't exist", func() {
				It("Should return an error", func() {
					By("Calling GetKey")
					mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)
					mockKvClient.EXPECT().
						Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

					client := kvclient.NewWithClient(mockKvClient)

					// TODO: implement specific error types in future by handling the gRPC error types.
					_, err := client.GetKey("test-collection", "test-key")

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		When("The collection doesn't exist", func() {
			It("Should return an error", func() {
				By("Calling GetKey")
				mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)
				mockKvClient.EXPECT().
					Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

				client := kvclient.NewWithClient(mockKvClient)

				// TODO: implement specific error types in future by handling the gRPC error types.
				_, err := client.GetKey("missing-collection", "test-key")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("DecodeKey", func() {

		When("The collection exists", func() {

			When("The key exists", func() {

				When("The key property names match the struct", func() {

					When("The key property types match the struct", func() {
						It("Should decode the keyvalue", func() {
							type DecodeKV struct {
								IntValue    int
								FloatValue  float32
								StringValue string
							}

							mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

							theValue := &DecodeKV{
								IntValue:    12,
								FloatValue:  12.5,
								StringValue: "test",
							}

							mockDocument, _ := structpb.NewStruct(structs.Map(theValue))

							By("Calling GetDocument")
							mockKvClient.EXPECT().
								Get(gomock.Any(), gomock.Any()).Return(&v1.KeyValueGetResponse{
								Value: mockDocument,
							}, nil)

							client := kvclient.NewWithClient(mockKvClient)
							theDecodedKv := DecodeKV{}

							err := client.DecodeKey("test-collection", "test-key", &theDecodedKv)

							By("Not returning an error")
							Expect(err).ShouldNot(HaveOccurred())

							By("Decoding the document")
							Expect(theDecodedKv.FloatValue).To(Equal(float32(12.5)))
							Expect(theDecodedKv.IntValue).To(Equal(12))
						})
					})

					When("The document property types don't match the struct", func() {
						It("Should return an error", func() {
							type Value struct {
								IntValue    int
								FloatValue  float32
								StringValue string
							}

							mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

							theValue := &Value{
								IntValue:    12,
								FloatValue:  12.5,
								StringValue: "test",
							}

							mockValue, _ := structpb.NewStruct(structs.Map(theValue))

							By("Calling DecodeKey")
							mockKvClient.EXPECT().
								Get(gomock.Any(), gomock.Any()).Return(&v1.KeyValueGetResponse{
								Value: mockValue,
							}, nil)

							client := kvclient.NewWithClient(mockKvClient)

							// Struct with the same property names, but incompatible types
							type TypeMismatchDocument struct {
								IntValue    string
								FloatValue  bool
								StringValue int
							}

							theDecodedDoc := TypeMismatchDocument{}
							err := client.DecodeKey("test-collection", "test-key", &theDecodedDoc)

							By("Returning an error")
							Expect(err).Should(HaveOccurred())
						})
					})
				})

				When("The struct contains an extra property", func() {
					It("Should decode the document, leaving the extra properties blank", func() {
						type Value struct {
							IntValue    int
							FloatValue  float32
							StringValue string
						}

						mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

						theValue := &Value{
							IntValue:    12,
							FloatValue:  12.5,
							StringValue: "test",
						}

						mockValue, _ := structpb.NewStruct(structs.Map(theValue))

						By("Calling GetKey")
						mockKvClient.EXPECT().
							Get(gomock.Any(), gomock.Any()).Return(&v1.KeyValueGetResponse{
							Value: mockValue,
						}, nil)

						client := kvclient.NewWithClient(mockKvClient)

						// Struct with the same property names, but incompatible types
						type TypeMismatchValue struct {
							IntValue    int
							FloatValue  float32
							StringValue string
							ExtraValue  string // One extra value, not present in the stored document.
						}

						theDecodedDoc := TypeMismatchValue{}
						err := client.DecodeKey("test-collection", "test-key", &theDecodedDoc)

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Leaving the extra property blank")
						Expect(theDecodedDoc.ExtraValue).To(Equal(""))
					})
				})

				When("The document contains an extra property", func() {
					It("Should return an error", func() {
						// Silently decoding documents into structs with keys missing could result in data loss
						// if those struct were subsequently used in a DocumentUpdate call.
						// It seems safest to allow structs to extend documents, but not be missing fields.

						type Value struct {
							IntValue    int
							FloatValue  float32
							StringValue string
							ExtraValue  string // One extra value, not present in the struct used for decoding.
						}

						mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

						theValue := &Value{
							IntValue:    12,
							FloatValue:  12.5,
							StringValue: "test",
							ExtraValue:  "extra",
						}

						mockStruct, _ := structpb.NewStruct(structs.Map(theValue))

						By("Calling GetKey")
						mockKvClient.EXPECT().
							Get(gomock.Any(), gomock.Any()).Return(&v1.KeyValueGetResponse{
							Value: mockStruct,
						}, nil)

						client := kvclient.NewWithClient(mockKvClient)

						// Struct with the same property names, but incompatible types
						type TypeMismatchValue struct {
							IntValue    int
							FloatValue  float32
							StringValue string
						}

						theDocodedValue := TypeMismatchValue{}
						err := client.DecodeKey("test-collection", "test-key", &theDocodedValue)

						By("Returning an error")
						Expect(err).Should(HaveOccurred())
					})

					When("Explicitly allowing unknown keys", func() {
						It("Should decode the document", func() {
							type Value struct {
								IntValue    int
								FloatValue  float32
								StringValue string
								ExtraValue  string // One extra value, not present in the struct used for decoding.
							}

							mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

							theDoc := &Value{
								IntValue:    12,
								FloatValue:  12.5,
								StringValue: "test",
								ExtraValue:  "extra",
							}

							mockValue, _ := structpb.NewStruct(structs.Map(theDoc))

							By("Calling GetDocument")
							mockKvClient.EXPECT().
								Get(gomock.Any(), gomock.Any()).Return(&v1.KeyValueGetResponse{
								Value: mockValue,
							}, nil)

							client := kvclient.NewWithClient(mockKvClient)

							// Struct with the same property names, but incompatible types
							type TypeMismatchValue struct {
								IntValue    int
								FloatValue  float32
								StringValue string
							}

							theDecodedDoc := TypeMismatchValue{}
							err := client.DecodeKey("test-collection", "test-key", &theDecodedDoc, kvclient.WithUnknownKeys(true))

							By("Returning an error")
							Expect(err).ShouldNot(HaveOccurred())
						})
					})
				})

				When("The document properties don't match the struct", func() {
					It("Should decode the document", func() {
						type DecodeValue struct {
							IntValue    int
							FloatValue  float32
							StringValue string
						}

						mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

						theValue := &DecodeValue{
							IntValue:    12,
							FloatValue:  12.5,
							StringValue: "test",
						}

						mockValue, _ := structpb.NewStruct(structs.Map(theValue))

						By("Calling GetDocument")
						mockKvClient.EXPECT().
							Get(gomock.Any(), gomock.Any()).Return(&v1.KeyValueGetResponse{
							Value: mockValue,
						}, nil)

						client := kvclient.NewWithClient(mockKvClient)
						theDecodedVal := DecodeValue{}

						err := client.DecodeKey("test-collection", "test-key", &theDecodedVal)

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Decoding the document")
						Expect(theDecodedVal.FloatValue).To(Equal(float32(12.5)))
						Expect(theDecodedVal.IntValue).To(Equal(12))
					})
				})
			})
		})
	})

	When("PutKey", func() {
		It("Should update the value", func() {
			mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

			By("Calling UpdateDocument with the expected inputs")
			mockKvClient.EXPECT().
				Put(
					gomock.Any(),
					&v1.KeyValuePutRequest{
						Collection: "test-collection",
						Key:        "test-key",
						Value: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"updated": structpb.NewNumberValue(321),
							},
						},
					},
				).Return(&v1.KeyValuePutResponse{}, nil)

			client := kvclient.NewWithClient(mockKvClient)
			input := map[string]interface{}{
				"updated": 321,
			}
			err := client.PutKey("test-collection", "test-key", input)

			By("Not returning an error")
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("DeleteKey", func() {
		When("The collection exists", func() {
			When("The document exists", func() {
				It("Should delete the key", func() {
					mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)

					By("Calling DeleteKey with the expected inputs")
					mockKvClient.EXPECT().
						Delete(
							gomock.Any(),
							&v1.KeyValueDeleteRequest{
								Collection: "test-collection",
								Key:        "test-key",
							},
						).Return(&v1.KeyValueDeleteResponse{}, nil)

					client := kvclient.NewWithClient(mockKvClient)
					err := client.DeleteKey("test-collection", "test-key")

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			When("The key doesn't exist", func() {
				It("Should return an error", func() {
					By("Calling DeleteKey")
					mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)
					mockKvClient.EXPECT().
						Delete(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

					client := kvclient.NewWithClient(mockKvClient)

					// TODO: implement specific error types in future by handling the gRPC error types.
					err := client.DeleteKey("test-collection", "test-key")

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		When("The collection doesn't exist", func() {
			It("Should return an error", func() {
				By("Calling DeleteDocument")
				mockKvClient := mock_v1.NewMockKeyValueClient(ctrl)
				mockKvClient.EXPECT().
					Delete(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

				client := kvclient.NewWithClient(mockKvClient)

				// TODO: implement specific error types in future by handling the gRPC error types.
				err := client.DeleteKey("test-collection", "test-key")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
