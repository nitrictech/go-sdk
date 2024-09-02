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

package topics

import (
	"context"

	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

// Topics
type Topics interface {
	// Topic - Retrieve a Topic reference
	Topic(name string) Topic
}

type topicsImpl struct {
	topicClient v1.TopicsClient
}

func (s *topicsImpl) Topic(name string) Topic {
	// Just return the straight topic reference
	// we can fail if the topic does not exist
	return &topicImpl{
		name:        name,
		topicClient: s.topicClient,
	}
}

// New - Construct a new Eventing Client with default options
func New() (Topics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.NitricDialTimeout())
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		return nil, errors.NewWithCause(codes.Unavailable, "Unable to dial Events service", err)
	}

	tc := v1.NewTopicsClient(conn)

	return &topicsImpl{
		topicClient: tc,
	}, nil
}
