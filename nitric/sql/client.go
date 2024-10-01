// Copyright 2023 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sql

import (
	"context"

	grpcx "github.com/nitrictech/go-sdk/internal/grpc"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"

	v1 "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
)

type SqlClientIface interface {
	// Name - The name of the store
	Name() string
	// Get a value from the store
	ConnectionString(context.Context) (string, error)
}

type SqlClient struct {
	name      string
	sqlClient v1.SqlClient
}

func (s *SqlClient) Name() string {
	return s.name
}

func (s *SqlClient) ConnectionString(ctx context.Context) (string, error) {
	resp, err := s.sqlClient.ConnectionString(ctx, &v1.SqlConnectionStringRequest{
		DatabaseName: s.name,
	})
	if err != nil {
		return "", err
	}

	return resp.ConnectionString, nil
}

func NewSqlClient(name string) (*SqlClient, error) {
	conn, err := grpcx.GetConnection()
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"unable to reach nitric server",
			err,
		)
	}

	client := v1.NewSqlClient(conn)

	return &SqlClient{
		name:      name,
		sqlClient: client,
	}, nil
}
