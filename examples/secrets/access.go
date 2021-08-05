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

package secrets_examples

// [START import]
import (
	"github.com/nitrictech/go-sdk/api/secrets"
)

// [END import]

func access() {
	// [START snippet]
	sc, _ := secrets.New()

	// Access the latest secret
	secVal, err := sc.Secret("my-secret").Latest().Access()

	if err != nil {
		// handle error
	}

	// do something with the secret
	_ = secVal.AsBytes()
	// [END snippet]
}
