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
	conn *grpc.ClientConn
	c v1.DocumentsClient
}

func (d NitricDocumentsClient) CreateDocument(collection string, key string, document map[string]interface{}) error  {
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

func (d NitricDocumentsClient) GetDocument(collection string, key string) (map[string]interface{}, error)  {
	res, err := d.c.GetDocument(context.Background(), &v1.GetDocumentRequest{
		Collection: collection,
		Key:        key,
	})
	if err != nil {
		return nil, err
	}
	return res.GetDocument().AsMap(), nil
}

func (d NitricDocumentsClient) UpdateDocument(collection string, key string, document map[string]interface{}) error  {
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

func (d NitricDocumentsClient) DeleteDocument(collection string, key string) error  {
	_, err := d.c.DeleteDocument(context.Background(), &v1.DeleteDocumentRequest{
		Collection: collection,
		Key:        key,
	})
	return err
}

// Close - closes the connection to the membrane server
// no need to call close if the connect is to remain open for the lifetime of the application.
func (d NitricDocumentsClient) Close() error {
	return d.conn.Close()
}

// FIXME: Extract into shared code.
func New() (DocumentsClient, error) {
	// Connect to the gRPC Membrane Server
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to Membrane gRPC server: %s", err)
	}

	return &NitricDocumentsClient{
		conn: conn,
		c: v1.NewDocumentsClient(conn),
	}, nil
}
