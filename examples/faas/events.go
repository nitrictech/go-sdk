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
