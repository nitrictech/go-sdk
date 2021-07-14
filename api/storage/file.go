// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	sc     v1.StorageServiceClient
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
