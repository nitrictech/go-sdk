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

	v1 "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
)

type Mode int

const (
	ModeRead Mode = iota
	ModeWrite
)

// File - A file reference for a bucket
type File interface {
	// Read - Read this object
	Read() ([]byte, error)
	// Write - Write this object
	Write([]byte) error
	// Delete - Delete this object
	Delete() error
	// PresignUrl - Creates a presigned Url for this file reference
	PresignUrl(PresignUrlOptions) (string, error)
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
		return nil, errors.FromGrpcError(err)
	}

	return r.GetBody(), nil
}

func (o *fileImpl) Write(content []byte) error {
	if _, err := o.sc.Write(context.TODO(), &v1.StorageWriteRequest{
		BucketName: o.bucket,
		Key:        o.key,
		Body:       content,
	}); err != nil {
		return errors.FromGrpcError(err)
	}

	return nil
}

func (o *fileImpl) Delete() error {
	if _, err := o.sc.Delete(context.TODO(), &v1.StorageDeleteRequest{
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

func (o *fileImpl) PresignUrl(opts PresignUrlOptions) (string, error) {
	if err := opts.isValid(); err != nil {
		return "", errors.NewWithCause(codes.InvalidArgument, "invalid options", err)
	}

	op := v1.StoragePreSignUrlRequest_READ

	if opts.Mode == ModeWrite {
		op = v1.StoragePreSignUrlRequest_WRITE
	}

	r, err := o.sc.PreSignUrl(context.TODO(), &v1.StoragePreSignUrlRequest{
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
