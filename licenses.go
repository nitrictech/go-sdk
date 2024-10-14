// Copyright 2021 Nitric Pty Ltd.
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

package main

// NOTE:
// This main package is a workaround for binary license scanning that forces transitive dependencies in
// Code we're distributing to be analyzed
import (
	_ "github.com/nitrictech/go-sdk/nitric"
	_ "github.com/nitrictech/go-sdk/nitric/apis"
	_ "github.com/nitrictech/go-sdk/nitric/batch"
	_ "github.com/nitrictech/go-sdk/nitric/errors"
	_ "github.com/nitrictech/go-sdk/nitric/keyvalue"
	_ "github.com/nitrictech/go-sdk/nitric/queues"
	_ "github.com/nitrictech/go-sdk/nitric/schedules"
	_ "github.com/nitrictech/go-sdk/nitric/secrets"
	_ "github.com/nitrictech/go-sdk/nitric/storage"
	_ "github.com/nitrictech/go-sdk/nitric/topics"
	_ "github.com/nitrictech/go-sdk/nitric/websockets"
)

func main() {}
