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

package nitric

import (
	"github.com/nitrictech/go-sdk/nitric/apis"
	"github.com/nitrictech/go-sdk/nitric/keyvalue"
	"github.com/nitrictech/go-sdk/nitric/queues"
	"github.com/nitrictech/go-sdk/nitric/schedules"
	"github.com/nitrictech/go-sdk/nitric/secrets"
	"github.com/nitrictech/go-sdk/nitric/sql"
	"github.com/nitrictech/go-sdk/nitric/storage"
	"github.com/nitrictech/go-sdk/nitric/topics"
	"github.com/nitrictech/go-sdk/nitric/websockets"
	"github.com/nitrictech/go-sdk/nitric/workers"
)

var (
	NewApi         = apis.NewApi
	NewKv          = keyvalue.NewKv
	NewQueue       = queues.NewQueue
	NewSchedule    = schedules.NewSchedule
	NewSecret      = secrets.NewSecret
	NewSqlDatabase = sql.NewSqlDatabase
	NewBucket      = storage.NewBucket
	NewTopic       = topics.NewTopic
	NewWebsocket   = websockets.NewWebsocket
)

func Run() error {
	return workers.GetDefaultManager().Run()
}
