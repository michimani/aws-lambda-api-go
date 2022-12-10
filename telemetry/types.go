package telemetry

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type DestinationProtocol string

const (
	DestinationProtocolHTTP DestinationProtocol = "HTTP"
	DestinationProtocolTCP  DestinationProtocol = "TCP"
)

func (dp DestinationProtocol) Valid() bool {
	return dp == DestinationProtocolHTTP || dp == DestinationProtocolTCP
}

type TelemetryType string

const (
	TelemetryTypePlatform  TelemetryType = "platform"
	TelemetryTypeFunction  TelemetryType = "function"
	TelemetryTypeExtension TelemetryType = "extension"
)

func (tt TelemetryType) Valid() bool {
	return tt == TelemetryTypePlatform || tt == TelemetryTypeFunction || tt == TelemetryTypeExtension
}

type SubscribeInput struct {
	// Generated unique identifier for public extension name.
	// This value will be got in response header of POST /extension/register API.
	LambdaExtensionIdentifier string

	// The protocol that Lambda uses to send telemetry data. (Required)
	DestinationProtocol DestinationProtocol

	// The URI to send telemetry data to. (Required)
	DestinationURI string

	// The types of telemetry that you want the extension to subscribe to. (Required)
	TelemetryTypes []TelemetryType

	// The maximum number of events to buffer in memory.
	// min/default/max = 25/1,000/30,000
	BufferMaxItems *uint64

	// The maximum volume of telemetry (in bytes) to buffer in memory.
	// min/default/max = 262,144/262,144/1,048,576
	BufferMaxBytes *uint64

	// The maximum time (in milliseconds) to buffer a batch.
	// min/default/max = 1,000/10,000/10,000
	BufferTimeoutMs *uint64
}

const schemaVersion = "2022-07-01"

type subscribeBody struct {
	SchemaVersion string                   `json:"schemaVersion"`
	Destination   subscribeBodyDestination `json:"destination"`
	Types         []string                 `json:"types"`
	Buffering     subscribeBodyBuffering   `json:"buffering"`
}

type subscribeBodyDestination struct {
	Protocol string `json:"protocol"`
	URI      string `json:"URI"`
}

type subscribeBodyBuffering struct {
	MaxItems  uint64 `json:"maxItems"`
	MaxBytes  uint64 `json:"maxBytes"`
	TimeoutMs uint64 `json:"timeoutMs"`
}

const (
	bufferingMaxItemsDefault  uint64 = 1000
	bufferingMaxBytesDefault  uint64 = 256 * 1024
	bufferingTimeoutMsDefault uint64 = 10000
)

var defaultBuffering = subscribeBodyBuffering{
	MaxItems:  bufferingMaxItemsDefault,
	MaxBytes:  bufferingMaxBytesDefault,
	TimeoutMs: bufferingTimeoutMsDefault,
}

func inputToRequestBody(in *SubscribeInput) (io.Reader, error) {
	if in == nil {
		return nil, errors.New("SubscribeInput is nil")
	}

	sb := subscribeBody{SchemaVersion: schemaVersion}

	if !in.DestinationProtocol.Valid() {
		return nil, errors.New("Invalid value for DestinationProtocol")
	}
	if len(in.DestinationURI) == 0 {
		return nil, errors.New("DestinationURI is empty")
	}
	sb.Destination = subscribeBodyDestination{
		Protocol: string(in.DestinationProtocol),
		URI:      in.DestinationURI,
	}

	if len(in.TelemetryTypes) == 0 {
		return nil, errors.New("TelemetryTypes is empty")
	}

	sbts := []string{}
	for _, tt := range in.TelemetryTypes {
		if !tt.Valid() {
			return nil, errors.New("TelemetryType has some invalid value")
		}
		sbts = append(sbts, string(tt))
	}
	sb.Types = sbts

	sb.Buffering = defaultBuffering
	if in.BufferMaxItems != nil {
		sb.Buffering.MaxItems = *in.BufferMaxItems
	}
	if in.BufferMaxBytes != nil {
		sb.Buffering.MaxBytes = *in.BufferMaxBytes
	}
	if in.BufferTimeoutMs != nil {
		sb.Buffering.TimeoutMs = *in.BufferTimeoutMs
	}

	j, err := json.Marshal(sb)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(j), nil
}

type SubscribeOutput struct {
	// http status code
	StatusCode int `json:"-"`

	// The error response.
	Error *ErrorResponse
}
type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorType    string `json:"errorType"`
}
