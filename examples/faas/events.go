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

package faas_examples

// [START import]
import (
	"encoding/json"
	"fmt"

	"github.com/nitrictech/go-sdk/api/events"
	"github.com/nitrictech/go-sdk/faas"
)

// [END import]

func evts() {
	// [START snippet]
	faas.New().Event(func(ctx *faas.EventContext, next faas.EventHandler) (*faas.EventContext, error) {
		var evt events.Event

		// Unmarshal a nitric event from an event payload (assuming it was published with the nitric events api)
		if err := json.Unmarshal(ctx.Request.Data(), &evt); err != nil {
			// Failed to handle the event
			ctx.Response.Success = false
			fmt.Println(err.Error())
		}

		fmt.Printf("Received nitric event: %v\n", evt.Payload)

		return ctx, nil
	}).Start()
	// [END snippet]
}
