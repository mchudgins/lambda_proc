package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/mchudgins/lambda_proc"
)

// Something to play with for the response

// TraceInfo provides the submitted input and a buffer to place log messages.
// This should be a normal part of your debugging workflow and included
// in the response json.
type TraceInfo struct {
	Log            string                 `json:"log,omitempty"`
	SubmittedInput map[string]interface{} `json:"submittedInput,omitempty"`
}

// Something is the Lambda API's response struct
type Something struct {
	Hello string     `json:"hello"`
	World string     `json:"world"`
	Trace *TraceInfo `json:"trace,omitempty"`
}

// APIEvent is a mapping of the Lambda event data from an API Gateway invocation.
// In this specific case, we use the API Gateway's Integration "mapping template"
// to map to a specific structure, rather than a map[string]interface{}
type APIEvent struct {
	Body        interface{}       `json:"body,omitempty"`
	Method      string            `json:"method"`
	Stage       string            `json:"stage"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers,omitempty"`
	QueryParams map[string]string `json:"query,omitempty"`
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

		var fTrace bool
		var v map[string]interface{}

		if err := json.Unmarshal(eventJSON, &v); err != nil {
			return nil, err
		}

		var evt APIEvent
		if err := json.Unmarshal(eventJSON, &evt); err != nil {
			return nil, err
		}
		log.Printf("Event:  %+v\n", evt)

		// check the query params
		if val, ok := evt.QueryParams["trace"]; ok {
			val = strings.ToLower(val)
			log.Printf("trace query param: '%s'", val)
			if strings.Compare(val, "yes") == 0 ||
				strings.Compare(val, "true") == 0 ||
				strings.Compare(val, "on") == 0 ||
				strings.Compare(val, "") == 0 {
				fTrace = true
			}
		}

		// some boilerplate logging.  Move to Init(), rather
		// than every request?
		var buf bytes.Buffer
		log := log.New(&buf, "golang_sample", log.Lshortfile)

		log.Printf("Version: %s\n", version)
		log.Printf("Build Time: %s\n", buildTime)
		log.Printf("Builder:  %s\n", builder)
		log.Printf("go Version:  %s\n", goversion)

		// per request boilerplate logging.
		log.Printf("Context: %+v\n", context)
		log.Printf("eventJSON: %+v\n", v)
		for key, value := range v {
			log.Printf("key: %s = %+v\n", key, value)
		}

		// business logic goes here
		switch strings.ToLower(evt.URL) {
		case "/widgets":
			log.Printf("URL: /widgets\n")
			response, err := ProcessWidgets(context, evt)
			return response, err

		case "":
			log.Printf("URL: /\n")

		default:
			log.Printf("URL: %s (unexpected!)\n", evt.URL)
		}

		//		build the response
		s := Something{Hello: context.FunctionName, World: context.FunctionVersion}

		// if '?trace' in the query params is true, then
		// allow the log to be returned to the user.
		if fTrace {
			s.Trace = &TraceInfo{SubmittedInput: v, Log: buf.String()}
		}

		return s, nil

	})

}
