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

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

// Cloud storage bucket resource for large file storage.
type Bucket interface {
	// File - Get a file reference for in this bucket
	File(key string) File
	// Files - Get all file references for this bucket
	Files(ctx context.Context) ([]File, error)
	// Name - Get the name of the bucket
	Name() string
}

type bucketImpl struct {
	storageClient   v1.StorageServiceClient
	name string
}

func (b *bucketImpl) File(key string) File {
	return &fileImpl{
		storageClient:     b.storageClient,
		bucket: b.name,
		key:    key,
	}
}

func (b *bucketImpl) Files(ctx context.Context) ([]File, error) {
	resp, err := b.storageClient.ListFiles(ctx, &v1.StorageListFilesRequest{
		BucketName: b.name,
	})
	if err != nil {
		return nil, err
	}

	fileRefs := make([]File, 0)

	for _, f := range resp.Files {
		fileRefs = append(fileRefs, &fileImpl{
			storageClient:     b.storageClient,
			bucket: b.name,
			key:    f.Key,
		})
	}

	return fileRefs, nil
}

func (b *bucketImpl) Name() string {
	return b.name
}