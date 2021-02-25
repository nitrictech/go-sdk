package documentsclient_test

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/golang/mock/gomock"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	"github.com/nitrictech/go-sdk/v1/documentsclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ = Describe("Documentsclient", func() {
	ctrl := gomock.NewController(GinkgoT())

	When("CreateDocument", func() {
		When("The collection exists", func() {
			It("Should create the document", func() {
				mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

				By("Calling CreateDocument with the expected inputs")
				mockDocClient.EXPECT().Create(gomock.Any(), &v1.DocumentCreateRequest{
					Collection: "test-collection",
					Key:        "test-key",
					Document:   &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"hello": structpb.NewNumberValue(123),
						},
					},
				})

				client := documentsclient.NewWithClient(mockDocClient)
				input := map[string]interface{}{
					"hello": 123,
				}
				err := client.CreateDocument("test-collection", "test-key", input)

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("The collection doesn't exist", func() {
			It("Return an error", func() {
				mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

				By("Calling CreateDocument with the expected inputs")
				mockDocClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

				client := documentsclient.NewWithClient(mockDocClient)
				input := map[string]interface{}{}
				err := client.CreateDocument("test-collection", "test-key", input)

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("GetDocument", func() {
		When("The collection exists", func() {
			When("The document exists", func() {
				It("Should retrieve the document", func() {
					mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

					mockDocument, _ := structpb.NewStruct(map[string]interface{}{
						"test": 123,
					})

					By("Calling GetDocument with the expected inputs")
					mockDocClient.EXPECT().
						Get(
							gomock.Any(),
							&v1.DocumentGetRequest{
								Collection: "test-collection",
								Key:        "test-key",
							},
						).Return(&v1.DocumentGetResponse{
						Document: mockDocument,
					}, nil)

					client := documentsclient.NewWithClient(mockDocClient)
					doc, err := client.GetDocument("test-collection", "test-key")

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the document")
					Expect(int(doc["test"].(float64))).To(Equal(123))
				})
			})
			When("The document doesn't exist", func() {
				It("Should return an error", func() {
					By("Calling GetDocument")
					mockDocClient := mock_v1.NewMockDocumentClient(ctrl)
					mockDocClient.EXPECT().
						Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

					client := documentsclient.NewWithClient(mockDocClient)

					// TODO: implement specific error types in future by handling the gRPC error types.
					_, err := client.GetDocument("test-collection", "test-key")

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		When("The collection doesn't exist", func() {
			It("Should return an error", func() {
				By("Calling GetDocument")
				mockDocClient := mock_v1.NewMockDocumentClient(ctrl)
				mockDocClient.EXPECT().
					Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

				client := documentsclient.NewWithClient(mockDocClient)

				// TODO: implement specific error types in future by handling the gRPC error types.
				_, err := client.GetDocument("missing-collection", "test-key")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("DecodeDocument", func() {

		When("The collection exists", func() {

			When("The document exists", func() {

				When("The document property names match the struct", func() {

					When("The document property types match the struct", func() {
						It("Should decode the document", func() {
							type DecodeDocument struct {
								IntValue    int
								FloatValue  float32
								StringValue string
							}

							mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

							theDoc := &DecodeDocument{
								IntValue:    12,
								FloatValue:  12.5,
								StringValue: "test",
							}

							mockDocument, _ := structpb.NewStruct(structs.Map(theDoc))

							By("Calling GetDocument")
							mockDocClient.EXPECT().
								Get(gomock.Any(), gomock.Any()).Return(&v1.DocumentGetResponse{
								Document: mockDocument,
							}, nil)

							client := documentsclient.NewWithClient(mockDocClient)
							theDecodedDoc := DecodeDocument{}

							err := client.DecodeDocument("test-collection", "test-key", &theDecodedDoc)

							By("Not returning an error")
							Expect(err).ShouldNot(HaveOccurred())

							By("Decoding the document")
							Expect(theDecodedDoc.FloatValue).To(Equal(float32(12.5)))
							Expect(theDecodedDoc.IntValue).To(Equal(12))
						})
					})

					When("The document property types don't match the struct", func() {
						It("Should return an error", func() {
							type Document struct {
								IntValue    int
								FloatValue  float32
								StringValue string
							}

							mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

							theDoc := &Document{
								IntValue:    12,
								FloatValue:  12.5,
								StringValue: "test",
							}

							mockDocument, _ := structpb.NewStruct(structs.Map(theDoc))

							By("Calling GetDocument")
							mockDocClient.EXPECT().
								Get(gomock.Any(), gomock.Any()).Return(&v1.DocumentGetResponse{
								Document: mockDocument,
							}, nil)

							client := documentsclient.NewWithClient(mockDocClient)

							// Struct with the same property names, but incompatible types
							type TypeMismatchDocument struct {
								IntValue    string
								FloatValue  bool
								StringValue int
							}

							theDecodedDoc := TypeMismatchDocument{}
							err := client.DecodeDocument("test-collection", "test-key", &theDecodedDoc)

							By("Returning an error")
							Expect(err).Should(HaveOccurred())
						})
					})
				})

				When("The struct contains an extra property", func() {
					It("Should decode the document, leaving the extra properties blank", func() {
						type Document struct {
							IntValue    int
							FloatValue  float32
							StringValue string
						}

						mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

						theDoc := &Document{
							IntValue:    12,
							FloatValue:  12.5,
							StringValue: "test",
						}

						mockDocument, _ := structpb.NewStruct(structs.Map(theDoc))

						By("Calling GetDocument")
						mockDocClient.EXPECT().
							Get(gomock.Any(), gomock.Any()).Return(&v1.DocumentGetResponse{
							Document: mockDocument,
						}, nil)

						client := documentsclient.NewWithClient(mockDocClient)

						// Struct with the same property names, but incompatible types
						type TypeMismatchDocument struct {
							IntValue    int
							FloatValue  float32
							StringValue string
							ExtraValue  string // One extra value, not present in the stored document.
						}

						theDecodedDoc := TypeMismatchDocument{}
						err := client.DecodeDocument("test-collection", "test-key", &theDecodedDoc)

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

						type Document struct {
							IntValue    int
							FloatValue  float32
							StringValue string
							ExtraValue  string // One extra value, not present in the struct used for decoding.
						}

						mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

						theDoc := &Document{
							IntValue:    12,
							FloatValue:  12.5,
							StringValue: "test",
							ExtraValue:  "extra",
						}

						mockDocument, _ := structpb.NewStruct(structs.Map(theDoc))

						By("Calling GetDocument")
						mockDocClient.EXPECT().
							Get(gomock.Any(), gomock.Any()).Return(&v1.DocumentGetResponse{
							Document: mockDocument,
						}, nil)

						client := documentsclient.NewWithClient(mockDocClient)

						// Struct with the same property names, but incompatible types
						type TypeMismatchDocument struct {
							IntValue    int
							FloatValue  float32
							StringValue string
						}

						theDecodedDoc := TypeMismatchDocument{}
						err := client.DecodeDocument("test-collection", "test-key", &theDecodedDoc)

						By("Returning an error")
						Expect(err).Should(HaveOccurred())
					})

					When("Explicitly allowing unknown keys", func() {
						It("Should decode the document", func() {
							type Document struct {
								IntValue    int
								FloatValue  float32
								StringValue string
								ExtraValue  string // One extra value, not present in the struct used for decoding.
							}

							mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

							theDoc := &Document{
								IntValue:    12,
								FloatValue:  12.5,
								StringValue: "test",
								ExtraValue:  "extra",
							}

							mockDocument, _ := structpb.NewStruct(structs.Map(theDoc))

							By("Calling GetDocument")
							mockDocClient.EXPECT().
								Get(gomock.Any(), gomock.Any()).Return(&v1.DocumentGetResponse{
								Document: mockDocument,
							}, nil)

							client := documentsclient.NewWithClient(mockDocClient)

							// Struct with the same property names, but incompatible types
							type TypeMismatchDocument struct {
								IntValue    int
								FloatValue  float32
								StringValue string
							}

							theDecodedDoc := TypeMismatchDocument{}
							err := client.DecodeDocument("test-collection", "test-key", &theDecodedDoc, documentsclient.WithUnknownKeys(true))

							By("Returning an error")
							Expect(err).ShouldNot(HaveOccurred())
						})
					})
				})

				When("The document properties don't match the struct", func() {
					It("Should decode the document", func() {
						type DecodeDocument struct {
							IntValue    int
							FloatValue  float32
							StringValue string
						}

						mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

						theDoc := &DecodeDocument{
							IntValue:    12,
							FloatValue:  12.5,
							StringValue: "test",
						}

						mockDocument, _ := structpb.NewStruct(structs.Map(theDoc))

						By("Calling GetDocument")
						mockDocClient.EXPECT().
							Get(gomock.Any(), gomock.Any()).Return(&v1.DocumentGetResponse{
							Document: mockDocument,
						}, nil)

						client := documentsclient.NewWithClient(mockDocClient)
						theDecodedDoc := DecodeDocument{}

						err := client.DecodeDocument("test-collection", "test-key", &theDecodedDoc)

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Decoding the document")
						Expect(theDecodedDoc.FloatValue).To(Equal(float32(12.5)))
						Expect(theDecodedDoc.IntValue).To(Equal(12))
					})
				})
			})
		})
	})

	When("UpdateDocument", func() {
		When("The collection exists", func() {
			When("The document exists", func() {
				It("Should update the document", func() {
					mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

					By("Calling UpdateDocument with the expected inputs")
					mockDocClient.EXPECT().
						Update(
							gomock.Any(),
							&v1.DocumentUpdateRequest{
								Collection: "test-collection",
								Key:        "test-key",
								Document:  &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"updated": structpb.NewNumberValue(321),
									},
								},
							},
						).Return(&v1.DocumentUpdateResponse{}, nil)

					client := documentsclient.NewWithClient(mockDocClient)
					input := map[string]interface{}{
						"updated": 321,
					}
					err := client.UpdateDocument("test-collection", "test-key", input)

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			When("The document doesn't exist", func() {
				It("Should return an error", func() {
					By("Calling UpdateDocument")
					mockDocClient := mock_v1.NewMockDocumentClient(ctrl)
					mockDocClient.EXPECT().
						Update(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

					client := documentsclient.NewWithClient(mockDocClient)

					// TODO: implement specific error types in future by handling the gRPC error types.
					err := client.UpdateDocument("missing-collection", "test-key", map[string]interface{}{})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		When("The collection doesn't exist", func() {
			It("Should return an error", func() {
				By("Calling UpdateDocument")
				mockDocClient := mock_v1.NewMockDocumentClient(ctrl)
				mockDocClient.EXPECT().
					Update(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

				client := documentsclient.NewWithClient(mockDocClient)

				// TODO: implement specific error types in future by handling the gRPC error types.
				err := client.UpdateDocument("missing-collection", "test-key", map[string]interface{}{})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("DeleteDocument", func() {
		When("The collection exists", func() {
			When("The document exists", func() {
				It("Should delete the document", func() {
					mockDocClient := mock_v1.NewMockDocumentClient(ctrl)

					By("Calling DeleteDocument with the expected inputs")
					mockDocClient.EXPECT().
						Delete(
							gomock.Any(),
							&v1.DocumentDeleteRequest{
								Collection: "test-collection",
								Key:        "test-key",
							},
						).Return(&v1.DocumentDeleteResponse{}, nil)

					client := documentsclient.NewWithClient(mockDocClient)
					err := client.DeleteDocument("test-collection", "test-key")

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			When("The document doesn't exist", func() {
				It("Should return an error", func() {
					By("Calling DeleteDocument")
					mockDocClient := mock_v1.NewMockDocumentClient(ctrl)
					mockDocClient.EXPECT().
						Delete(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

					client := documentsclient.NewWithClient(mockDocClient)

					// TODO: implement specific error types in future by handling the gRPC error types.
					err := client.DeleteDocument("test-collection", "test-key")

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		When("The collection doesn't exist", func() {
			It("Should return an error", func() {
				By("Calling DeleteDocument")
				mockDocClient := mock_v1.NewMockDocumentClient(ctrl)
				mockDocClient.EXPECT().
					Delete(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock not found error"))

				client := documentsclient.NewWithClient(mockDocClient)

				// TODO: implement specific error types in future by handling the gRPC error types.
				err := client.DeleteDocument("test-collection", "test-key")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
