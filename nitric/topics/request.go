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

package topics

type Request interface {
	TopicName() string
	Message() map[string]interface{}
}

type requestImpl struct {
	topicName string
	message   map[string]interface{}
}

func (m *requestImpl) TopicName() string {
	return m.topicName
}

func (m *requestImpl) Message() map[string]interface{} {
	return m.message
}
