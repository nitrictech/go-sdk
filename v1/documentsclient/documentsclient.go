package documentsclient

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type DocumentsClient interface {
	CreateDocument(collection string, key string, document map[string]interface{}) error
	GetDocument(collection string, key string) (map[string]interface{}, error)
	DecodeDocument(collection string, key string, output interface{}, opts ...DecodeOption) error
	UpdateDocument(collection string, key string, document map[string]interface{}) error
	DeleteDocument(collection string, key string) error
}

type NitricDocumentsClient struct {
	c v1.DocumentClient
}

// CreateDocument - stores a new document in the document db
func (d NitricDocumentsClient) CreateDocument(collection string, key string, document map[string]interface{}) error {
	// Convert payload to Protobuf Struct
	docStruct, err := structpb.NewStruct(document)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %s", err)
	}

	_, err = d.c.Create(context.Background(), &v1.DocumentCreateRequest{
		Collection: collection,
		Key:        key,
		Document:   docStruct,
	})

	return err
}

type DecodeOption interface {
	Apply(c *mapstructure.DecoderConfig)
}

func WithUnknownKeys(allow bool) DecodeOption {
	return withUnknownKeys{allow}
}

type withUnknownKeys struct { allow bool }

func (w withUnknownKeys) Apply(c *mapstructure.DecoderConfig)  {
	c.ErrorUnused = !w.allow
}

// DecodeDocument - retrieves a document and decodes its contents into the given Go interface{}
//
// internally this method calls GetDocument then decodes the map[string]interface{} into the supplied interface{}
//
// this method helps parse the types of documents represented by structs.
func (d NitricDocumentsClient) DecodeDocument(collection string, key string, output interface{}, opts ...DecodeOption) error {
	document, err := d.GetDocument(collection, key)
	if err != nil {
		return err
	}
	decoderConfig := mapstructure.DecoderConfig{
		//DecodeHook:       nil,
		ErrorUnused:      true, // Default behavior is to error when keys are missing from the output interface{}
		//ZeroFields:       false,
		//WeaklyTypedInput: false,
		//Squash:           false,
		//Metadata:         nil,
		Result:           output,
		//TagName:          "",
	}

	// Apply additional options
	for _, opt := range opts {
		opt.Apply(&decoderConfig)
	}

	// Decode the document into the object
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return err
	}
	return decoder.Decode(document)
}

// GetDocument - retrieve an existing document from the document db
func (d NitricDocumentsClient) GetDocument(collection string, key string) (map[string]interface{}, error) {
	res, err := d.c.Get(context.Background(), &v1.DocumentGetRequest{
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

	_, err = d.c.Update(context.Background(), &v1.DocumentUpdateRequest{
		Collection: collection,
		Key:        key,
		Document:   docStruct,
	})

	return err
}

// DeleteDocument - deletes an existing document from the document db
func (d NitricDocumentsClient) DeleteDocument(collection string, key string) error {
	_, err := d.c.Delete(context.Background(), &v1.DocumentDeleteRequest{
		Collection: collection,
		Key:        key,
	})
	return err
}

func NewDocumentsClient(conn *grpc.ClientConn) DocumentsClient {
	return &NitricDocumentsClient{
		c: v1.NewDocumentClient(conn),
	}
}

func NewWithClient(client v1.DocumentClient) DocumentsClient {
	return &NitricDocumentsClient{
		c: client,
	}
}
