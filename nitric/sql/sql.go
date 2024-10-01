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

package sql

import (
	"github.com/nitrictech/go-sdk/nitric/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type sqlDatabaseOption func(*v1.SqlDatabaseResource)

func WithMigrationsPath(path string) sqlDatabaseOption {
	return func(r *v1.SqlDatabaseResource) {
		r.Migrations = &v1.SqlDatabaseMigrations{
			Migrations: &v1.SqlDatabaseMigrations_MigrationsPath{
				MigrationsPath: path,
			},
		}
	}
}

// NewSqlDatabase - Create a new Sql Database resource
func NewSqlDatabase(name string, opts ...sqlDatabaseOption) (*SqlClient, error) {
	resourceConfig := &v1.ResourceDeclareRequest_SqlDatabase{
		SqlDatabase: &v1.SqlDatabaseResource{},
	}

	for _, opt := range opts {
		opt(resourceConfig.SqlDatabase)
	}

	registerChan := workers.GetDefaultManager().RegisterResource(&v1.ResourceDeclareRequest{
		Id: &v1.ResourceIdentifier{
			Type: v1.ResourceType_SqlDatabase,
			Name: name,
		},
		Config: resourceConfig,
	})

	// Make sure that registerChan is read
	// Currently sql databases do not have allow methods so there is no reason to block on this
	go func() {
		<-registerChan
	}()

	return NewSqlClient(name)
}
