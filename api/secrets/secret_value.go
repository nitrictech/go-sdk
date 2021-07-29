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

package secrets

// SecretValue - Interface
type SecretValue interface {
	// SecretVersionRef - Returns a reference to the version of this value
	Ref() SecretVersionRef
	// AsBytes - Returns the secret value as bytes
	AsBytes() []byte
	// AsString - Returns the secret value as a string
	AsString() string
}

type secretValueImpl struct {
	version SecretVersionRef
	val     []byte
}

func (s *secretValueImpl) Ref() SecretVersionRef {
	return s.version
}

func (s *secretValueImpl) AsBytes() []byte {
	return s.val
}

func (s *secretValueImpl) AsString() string {
	return string(s.AsBytes())
}
