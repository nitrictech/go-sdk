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

package grpcx

import (
	"sync"

	"github.com/nitrictech/go-sdk/constants"
	"google.golang.org/grpc"
)

type grpcManager struct {
	conn      grpc.ClientConnInterface
	connMutex sync.Mutex
}

var m = grpcManager{
	conn:      nil,
	connMutex: sync.Mutex{},
}

func GetConnection() (grpc.ClientConnInterface, error) {
	m.connMutex.Lock()
	defer m.connMutex.Unlock()

	if m.conn == nil {
		conn, err := grpc.NewClient(constants.NitricAddress(), constants.DefaultOptions()...)
		if err != nil {
			return nil, err
		}
		m.conn = conn
	}

	return m.conn, nil
}
