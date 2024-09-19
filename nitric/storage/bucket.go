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
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

// Cloud storage bucket resource for large file storage.
type BucketIface interface {
	// Name - Get the name of the bucket
	Name() string
	// ListFiles - List the files in the bucket
	ListFiles(ctx context.Context) ([]string, error)
	// Read - Read this object
	Read(ctx context.Context, key string) ([]byte, error)
	// Write - Write this object
	Write(ctx context.Context, key string, data []byte) error
	// Delete - Delete this object
	Delete(ctx context.Context, key string) error
	// UploadUrl - Creates a signed Url for uploading this file reference
	UploadUrl(ctx context.Context, key string, opts ...PresignUrlOption) (string, error)
	// DownloadUrl - Creates a signed Url for downloading this file reference
	DownloadUrl(ctx context.Context, key string, opts ...PresignUrlOption) (string, error)
}

var _ BucketIface = (*Bucket)(nil)

type Bucket struct {
	storageClient v1.StorageClient
	name          string
}

func (o *Bucket) Read(ctx context.Context, key string) ([]byte, error) {
	r, err := o.storageClient.Read(ctx, &v1.StorageReadRequest{
		BucketName: o.name,
		Key:        key,
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	return r.GetBody(), nil
}

func (o *Bucket) Write(ctx context.Context, key string, content []byte) error {
	if _, err := o.storageClient.Write(ctx, &v1.StorageWriteRequest{
		BucketName: o.name,
		Key:        key,
		Body:       content,
	}); err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

func (o *Bucket) Delete(ctx context.Context, key string) error {
	if _, err := o.storageClient.Delete(ctx, &v1.StorageDeleteRequest{
		BucketName: o.name,
		Key:        key,
	}); err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

func (b *Bucket) ListFiles(ctx context.Context) ([]string, error) {
	resp, err := b.storageClient.ListBlobs(ctx, &v1.StorageListBlobsRequest{
		BucketName: b.name,
	})
	if err != nil {
		return nil, err
	}

	fileRefs := make([]string, 0)

	for _, f := range resp.Blobs {
		fileRefs = append(fileRefs, f.Key)
	}

	return fileRefs, nil
}

type Mode int

const (
	ModeRead Mode = iota
	ModeWrite
)

type presignUrlOptions struct {
	mode   Mode
	expiry time.Duration
}

type PresignUrlOption func(opts *presignUrlOptions)

func WithPresignUrlExpiry(expiry time.Duration) PresignUrlOption {
	return func(opts *presignUrlOptions) {
		opts.expiry = expiry
	}
}

func getPresignUrlOpts(mode Mode, opts ...PresignUrlOption) *presignUrlOptions {
	defaultOpts := &presignUrlOptions{
		mode:   mode,
		expiry: time.Minute * 5,
	}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	return defaultOpts
}

func (o *Bucket) UploadUrl(ctx context.Context, key string, opts ...PresignUrlOption) (string, error) {
	optsWithDefaults := getPresignUrlOpts(ModeWrite, opts...)

	return o.signUrl(ctx, key, optsWithDefaults)
}

func (o *Bucket) DownloadUrl(ctx context.Context, key string, opts ...PresignUrlOption) (string, error) {
	optsWithDefaults := getPresignUrlOpts(ModeRead, opts...)

	return o.signUrl(ctx, key, optsWithDefaults)
}

func (o *Bucket) signUrl(ctx context.Context, key string, opts *presignUrlOptions) (string, error) {
	op := v1.StoragePreSignUrlRequest_READ

	if opts.mode == ModeWrite {
		op = v1.StoragePreSignUrlRequest_WRITE
	}

	r, err := o.storageClient.PreSignUrl(ctx, &v1.StoragePreSignUrlRequest{
		BucketName: o.name,
		Key:        key,
		Operation:  op,
		Expiry:     durationpb.New(opts.expiry),
	})
	if err != nil {
		return "", errors.FromGrpcError(err)
	}

	return r.Url, nil
}

func (b *Bucket) Name() string {
	return b.name
}

func NewBucket(name string) (*Bucket, error) {
	conn, err := grpc.NewClient(constants.NitricAddress(), constants.DefaultOptions()...)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"unable to reach nitric server",
			err,
		)
	}

	storageClient := v1.NewStorageClient(conn)

	return &Bucket{
		name:          name,
		storageClient: storageClient,
	}, nil
}
