package extension

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type EventType string

const (
	EventTypeInvoke   EventType = "INVOKE"
	EventTypeShutdown EventType = "SHUTDOWN"
)

type events struct {
	Events []EventType `json:"events"`
}

func (es *events) toRequestBody() (io.Reader, error) {
	if es == nil {
		return nil, errors.New("es is nil")
	}

	j, err := json.Marshal(es)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(j), nil
}

type RegisterInput struct {
	// Public extension name.
	LambdaExtensionName string

	// Use this to specify optional Extensions features. Comma separated string. Available features:
	// * `accountId` - the register response will contain the account id associated with the Lambda function for which the Extension is being registered
	LambdaExtensionAcceptFeature string

	// EventTypes that the extension want to receive.
	Events []EventType
}

type RegisterOutput struct {
	// http status code
	StatusCode int `json:"-"`

	// Generated unique identifier for public extension name.
	LambdaExtensionIdentifier string `json:"-"`

	// Function Name.
	FunctionName string `json:"functionName"`

	// Function version. (e.g. $LATEST)
	FunctionVersion string `json:"functionVersion"`

	// Handler of the function.
	Handler string `json:"handler"`

	// AWS AccountID. Filled with only registration with `accountId` feature.
	AccountID string `json:"accountId"`

	// The error response.
	Error *ErrorResponse
}

type EventNextInput struct {
	// Unique identifier for extension.
	LambdaExtensionIdentifier string
}

type EventNextOutput struct {
	// http status code
	StatusCode int `json:"-"`

	// Value of Lambda-Extension-Event-Identifier header
	LambdaExtensionEventIdentifier string `json:"-"`

	// Type of next event. INVOKE | SHUTDOWN
	EventType string `json:"eventType"`

	// Function execution deadline counted in milliseconds since the Unix epoch.
	DeadlineMs int `json:"deadlineMs"`

	// AWS request ID associated with the request. Filled only with event type INVOKE.
	RequestID string `json:"requestId"`

	// The ARN requested. This can be different in each invoke that	executes the same version. Filled only with event type INVOKE.
	InvokedFunctionArn string `json:"invokedFunctionArn"`

	// X-Ray tracing value. Filled only with event type INVOKE.
	Tracing XRayTracingInfo `json:"tracing"`

	// Reason of shutdown. Filled only with event type SHUTDOWN.
	ShutdownReason string `json:"shutdownReason"`

	// The error response.
	Error *ErrorResponse `json:"-"`
}

type XRayTracingInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorType    string `json:"errorType"`
}
