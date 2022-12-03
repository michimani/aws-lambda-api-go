package extension

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
