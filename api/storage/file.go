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
	"fmt"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	v1 "github.com/nitrictech/go-sdk/nitric/v1"
)

type Mode int

const (
	ModeRead Mode = iota
	ModeWrite
)

// File - A file reference for a bucket
type File interface {
	// Name - Get the name of the file
	Name() string
	// Read - Read this object
	Read(ctx context.Context) ([]byte, error)
	// Write - Write this object
	Write(ctx context.Context, data []byte) error
	// Delete - Delete this object
	Delete(ctx context.Context) error
	// UploadUrl - Creates a signed Url for uploading this file reference
	UploadUrl(ctx context.Context, expiry int) (string, error)
	// DownloadUrl - Creates a signed Url for downloading this file reference
	DownloadUrl(ctx context.Context, expiry int) (string, error)
}

type fileImpl struct {
	bucket string
	key    string
	sc     v1.StorageServiceClient
}

func (o *fileImpl) Name() string {
	return o.key
}

func (o *fileImpl) Read(ctx context.Context) ([]byte, error) {
	r, err := o.sc.Read(ctx, &v1.StorageReadRequest{
		BucketName: o.bucket,
		Key:        o.key,
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	return r.GetBody(), nil
}

func (o *fileImpl) Write(ctx context.Context, content []byte) error {
	if _, err := o.sc.Write(ctx, &v1.StorageWriteRequest{
		BucketName: o.bucket,
		Key:        o.key,
		Body:       content,
	}); err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

func (o *fileImpl) Delete(ctx context.Context) error {
	if _, err := o.sc.Delete(ctx, &v1.StorageDeleteRequest{
		BucketName: o.bucket,
		Key:        o.key,
	}); err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

type PresignUrlOptions struct {
	Mode   Mode
	Expiry int
}

func (p PresignUrlOptions) isValid() error {
	if p.Mode != ModeRead && p.Mode != ModeWrite {
		return fmt.Errorf("invalid mode: %d", p.Mode)
	}

	return nil
}

func (o *fileImpl) UploadUrl(ctx context.Context, expiry int) (string, error) {
	return o.signUrl(ctx, PresignUrlOptions{Expiry: expiry, Mode: ModeWrite})
}

func (o *fileImpl) DownloadUrl(ctx context.Context, expiry int) (string, error) {
	return o.signUrl(ctx, PresignUrlOptions{Expiry: expiry, Mode: ModeRead})
}

func (o *fileImpl) signUrl(ctx context.Context, opts PresignUrlOptions) (string, error) {
	if err := opts.isValid(); err != nil {
		return "", errors.NewWithCause(codes.InvalidArgument, "invalid options", err)
	}

	op := v1.StoragePreSignUrlRequest_READ

	if opts.Mode == ModeWrite {
		op = v1.StoragePreSignUrlRequest_WRITE
	}

	r, err := o.sc.PreSignUrl(ctx, &v1.StoragePreSignUrlRequest{
		BucketName: o.bucket,
		Key:        o.key,
		Operation:  op,
		Expiry:     uint32(opts.Expiry),
	})
	if err != nil {
		return "", errors.FromGrpcError(err)
	}

	return r.Url, nil
}
