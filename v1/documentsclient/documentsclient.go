package documentsclient

import (
	"context"
	"fmt"

	v1 "go.nitric.io/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type DocumentsClient interface {
	CreateDocument(collection string, key string, document map[string]interface{}) error
	GetDocument(collection string, key string) (map[string]interface{}, error)
	UpdateDocument(collection string, key string, document map[string]interface{}) error
	DeleteDocument(collection string, key string) error
}

type NitricDocumentsClient struct {
	c v1.DocumentsClient
}

// CreateDocument - stores a new document in the document db
func (d NitricDocumentsClient) CreateDocument(collection string, key string, document map[string]interface{}) error {
	// Convert payload to Protobuf Struct
	docStruct, err := structpb.NewStruct(document)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %s", err)
	}

	_, err = d.c.CreateDocument(context.Background(), &v1.CreateDocumentRequest{
		Collection: collection,
		Key:        key,
		Document:   docStruct,
	})

	return err
}

// GetDocument - retrieve an existing document from the document db
func (d NitricDocumentsClient) GetDocument(collection string, key string) (map[string]interface{}, error) {
	res, err := d.c.GetDocument(context.Background(), &v1.GetDocumentRequest{
		Collection: collection,
		Key:        key,
	})
	if err != nil {
		return nil, err
	}
	return res.GetDocument().AsMap(), nil
}

// UpdateDocument - updates the contents of an existing document in the document db
func (d NitricDocumentsClient) UpdateDocument(collection string, key string, document map[string]interface{}) error {
	// Convert payload to Protobuf Struct
	docStruct, err := structpb.NewStruct(document)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %s", err)
	}

	_, err = d.c.UpdateDocument(context.Background(), &v1.UpdateDocumentRequest{
		Collection: collection,
		Key:        key,
		Document:   docStruct,
	})

	return err
}

// DeleteDocument - deletes an existing document from the document db
func (d NitricDocumentsClient) DeleteDocument(collection string, key string) error {
	_, err := d.c.DeleteDocument(context.Background(), &v1.DeleteDocumentRequest{
		Collection: collection,
		Key:        key,
	})
	return err
}

func NewDocumentsClient(conn *grpc.ClientConn) DocumentsClient {
	return &NitricDocumentsClient{
		c: v1.NewDocumentsClient(conn),
	}
}

func NewWithClient(client v1.DocumentsClient) DocumentsClient {
	return &NitricDocumentsClient{
		c: client,
	}
}
