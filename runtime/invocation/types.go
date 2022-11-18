package invocation

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/michimani/aws-lambda-api-go/runtime"
)

type NextOutput struct {
	// http status code
	StatusCode int

	// AWS request ID associated with the request.
	AWSRequestID string

	// X-Ray tracing header.
	TraceID string

	// Information about the client application and device when invoked	through the AWS Mobile SDK.
	ClientContext string

	// Information about the Amazon Cognito identity provider when invoked through the AWS Mobile SDK.
	CognitoIdentity string

	// Function execution deadline counted in milliseconds since the Unix epoch.
	DeadlineMs string

	// The ARN requested. This can be different in each invoke that	executes the same version.
	InvokedFunctionArn string

	// The bytes of EventResponse.
	RawEventResponse []byte

	// The error response.
	Error *runtime.ErrorResponse
}

// UnmarshalEventResponse converts the EventResponse returned as the response body of Runtime API
// into the structure received as an argument.
func (o *NextOutput) UnmarshalEventResponse(target any) error {
	if o == nil {
		return errors.New("Receiver is nil.")
	}

	if err := json.Unmarshal(o.RawEventResponse, target); err != nil {
		return err
	}

	return nil
}

// ResponseInput is the struct for parameter of
// POST /runtime/invocation/{AwsRequestId}/response API.
type ResponseInput struct {
	// AWS request ID associated with the request.
	AWSRequestID string

	// Any data that will be returned after the function has run to completion.
	Response io.Reader
}

// ResponseOutput is the struct for response of
// POST /runtime/invocation/{AwsRequestId}/response API.
type ResponseOutput struct {
	// http status code
	StatusCode int `json:"-"`

	// status
	Status string `json:"status"`

	// The error response
	Error *runtime.ErrorResponse `json:"-"`
}
