package storage

import v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"

type Bucket interface {
	// Object - Get an object reference for in this bucket
	File(key string) File
}

type bucketImpl struct {
	sc   v1.StorageClient
	name string
}

func (b *bucketImpl) File(key string) File {
	return &fileImpl{
		sc:     b.sc,
		bucket: b.name,
		key:    key,
	}
}
