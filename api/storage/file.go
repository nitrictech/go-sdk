package storage

import (
	"context"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
)

// File - A file reference for a bucket
type File interface {
	// Read - Read this object
	Read() ([]byte, error)
	// Write - Write this object
	Write([]byte) error
	// Delete - Delete this object
	Delete() error
}

type fileImpl struct {
	bucket string
	key    string
	sc     v1.StorageClient
}

func (o *fileImpl) Read() ([]byte, error) {
	r, err := o.sc.Read(context.TODO(), &v1.StorageReadRequest{
		BucketName: o.bucket,
		Key:        o.key,
	})

	if err != nil {
		return nil, err
	}

	return r.GetBody(), nil
}

func (o *fileImpl) Write(content []byte) error {
	_, err := o.sc.Write(context.TODO(), &v1.StorageWriteRequest{
		BucketName: o.bucket,
		Key:        o.key,
		Body:       content,
	})

	return err
}

func (o *fileImpl) Delete() error {
	_, err := o.sc.Delete(context.TODO(), &v1.StorageDeleteRequest{
		BucketName: o.bucket,
		Key:        o.key,
	})

	return err
}
