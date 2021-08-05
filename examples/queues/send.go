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

package queues_examples

// [START import]
import (
	"fmt"

	"github.com/nitrictech/go-sdk/api/queues"
)

// [END import]

func send() {
	// [START snippet]
	qc, _ := queues.New()

	ftasks, err := qc.Queue("my-queue").Send([]*queues.Task{
		{
			ID:          "1234",
			PayloadType: "test-payload",
			Payload: map[string]interface{}{
				"example": "payload",
			},
		},
	})

	if err != nil {
		// handle error...
	}

	for _, ft := range ftasks {
		// enumerate tasks that failed to queue
		fmt.Println("Failed Task", ft)
	}
	// [END snippet]
}
