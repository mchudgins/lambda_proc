package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/mchudgins/lambda_proc"
)

// Something to play with for the response

// TraceInfo provides the submitted input and a buffer to place log messages.
// This should be a normal part of your debugging workflow and included
// in the response json.
type TraceInfo struct {
	Log            string                 `json:"log"`
	SubmittedInput map[string]interface{} `json:"submittedInput"`
}

// Something is the Lambda API's response struct
type Something struct {
	Hello string    `json:"hello"`
	World string    `json:"world"`
	Trace TraceInfo `json:"trace"`
}

var (
	// boilerplate variables for good SDLC hygiene.  These are auto-magically
	// injected by the Makefile & linker working together.
	version   string
	buildTime string
	builder   string
	goversion string
)

func main() {

	// Define and run the Lambda function
	lambda_proc.Run(func(context *lambda_proc.Context, eventJSON json.RawMessage) (interface{}, error) {

		var v map[string]interface{}

		if err := json.Unmarshal(eventJSON, &v); err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		log := log.New(&buf, "golang_sample", log.Lshortfile)

		// if something in the query params is true, then
		log.Printf("Version: %s\n", version)
		log.Printf("Build Time: %s\n", buildTime)
		log.Printf("Builder:  %s\n", builder)
		log.Printf("go Version:  %s\n", goversion)
		log.Printf("Context: %+v\n", context)
		log.Printf("eventJSON: %+v\n", v)
		for key, value := range v {
			log.Printf("key: %s = %+v\n", key, value)
		}

		// business logic goes here

		//		build the response
		s := Something{Hello: "world", World: "latest", Trace: TraceInfo{SubmittedInput: v, Log: buf.String()}}

		return s, nil

	})

}
