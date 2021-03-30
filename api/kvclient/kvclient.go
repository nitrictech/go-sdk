package kvclient

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type KVClient interface {
	GetKey(collection string, key string) (map[string]interface{}, error)
	DecodeKey(collection string, key string, output interface{}, opts ...DecodeOption) error
	PutKey(collection string, key string, value map[string]interface{}) error
	DeleteKey(collection string, key string) error
}

type NitricKVClient struct {
	c v1.KeyValueClient
}

type DecodeOption interface {
	Apply(c *mapstructure.DecoderConfig)
}

func WithUnknownKeys(allow bool) DecodeOption {
	return withUnknownKeys{allow}
}

type withUnknownKeys struct{ allow bool }

func (w withUnknownKeys) Apply(c *mapstructure.DecoderConfig) {
	c.ErrorUnused = !w.allow
}

// DecodeKey - retrieves a value and decodes its contents into the given Go interface{}
//
// internally this method calls GetKey then decodes the map[string]interface{} into the supplied interface{}
//
// this method helps parse the types of value represented by structs.
func (d NitricKVClient) DecodeKey(collection string, key string, output interface{}, opts ...DecodeOption) error {
	value, err := d.GetKey(collection, key)
	if err != nil {
		return err
	}
	decoderConfig := mapstructure.DecoderConfig{
		//DecodeHook:       nil,
		ErrorUnused: true, // Default behavior is to error when keys are missing from the output interface{}
		//ZeroFields:       false,
		//WeaklyTypedInput: false,
		//Squash:           false,
		//Metadata:         nil,
		Result: output,
		//TagName:          "",
	}

	// Apply additional options
	for _, opt := range opts {
		opt.Apply(&decoderConfig)
	}

	// Decode the value into the object
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return err
	}
	return decoder.Decode(value)
}

// GetKey - retrieve an existing value from the kv store
func (d NitricKVClient) GetKey(collection string, key string) (map[string]interface{}, error) {
	res, err := d.c.Get(context.Background(), &v1.KeyValueGetRequest{
		Collection: collection,
		Key:        key,
	})
	if err != nil {
		return nil, err
	}
	return res.GetValue().AsMap(), nil
}

// PutKey - updates the value of an existing key in the kv store
func (d NitricKVClient) PutKey(collection string, key string, value map[string]interface{}) error {
	// Convert payload to Protobuf Struct
	valueStruct, err := structpb.NewStruct(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value: %s", err)
	}

	_, err = d.c.Put(context.Background(), &v1.KeyValuePutRequest{
		Collection: collection,
		Key:        key,
		Value:      valueStruct,
	})

	return err
}

// DeleteKey - deletes an existing key from the kv store
func (d NitricKVClient) DeleteKey(collection string, key string) error {
	_, err := d.c.Delete(context.Background(), &v1.KeyValueDeleteRequest{
		Collection: collection,
		Key:        key,
	})
	return err
}

func NewKVClient(conn *grpc.ClientConn) KVClient {
	return &NitricKVClient{
		c: v1.NewKeyValueClient(conn),
	}
}

func NewWithClient(client v1.KeyValueClient) KVClient {
	return &NitricKVClient{
		c: client,
	}
}
