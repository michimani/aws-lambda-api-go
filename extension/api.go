package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/internal"
)

const (
	registerEndpointFmt  string = "http://%s/2020-01-01/extension/register"
	eventNextEndpointFmt string = "http://%s/2020-01-01/extension/event/next"

	// Request Header
	requestHeaderLambdaExtensionName          string = "Lambda-Extension-Name"
	requestHeaderLambdaExtensionAcceptFeature string = "Lambda-Extension-Accept-Feature"
	requestHeaderLambdaExtensionIdentifier    string = "Lambda-Extension-Identifier"

	// Response Header
	responseHeaderLambdaExtensionEventIdentifier string = "Lambda-Extension-Event-Identifier"
	responseHeaderLambdaExtensionIdentifier      string = "Lambda-Extension-Identifier"
)

// Register an extension with the given name.
//
// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-extensions-api.html#extensions-api-next
func Register(ctx context.Context, client alago.AlagoClient, in *RegisterInput) (*RegisterOutput, error) {
	if in == nil {
		return nil, fmt.Errorf("RegisterInput is nil")
	}
	if in.LambdaExtensionName == "" {
		return nil, fmt.Errorf("RegisterInput.LambdaExtensionName is empty")
	}
	hs := []internal.Header{
		{Key: requestHeaderLambdaExtensionName, Value: in.LambdaExtensionName},
	}

	if in.LambdaExtensionAcceptFeature == "" {
		hs = append(hs, internal.Header{Key: requestHeaderLambdaExtensionAcceptFeature, Value: in.LambdaExtensionAcceptFeature})
	}

	ev := events{Events: in.Events}
	reqBody, err := ev.toRequestBody()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(registerEndpointFmt, client.Host())
	sc, h, b, err := internal.CallAPI(context.Background(), client, http.MethodPost, url, reqBody, hs...)
	if err != nil {
		return nil, err
	}

	out, err := generateRegisterOutput(sc, h, b)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func generateRegisterOutput(sc int, header http.Header, body []byte) (*RegisterOutput, error) {
	out := RegisterOutput{}
	out.StatusCode = sc

	if v, ok := header[responseHeaderLambdaExtensionIdentifier]; ok && len(v) > 0 {
		out.LambdaExtensionIdentifier = v[0]
	}

	if sc != http.StatusOK {
		var errRes ErrorResponse
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}
		out.Error = &errRes
		return &out, nil
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("err:%v, body:%s", err, string(body))
	}

	return &out, nil
}

// Extension makes this HTTP request when it is ready to receive and process a new event.
// This is an iterator-style blocking API call. Response contains event JSON document.
//
// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-extensions-api.html#extensions-api-next
func EventNext(ctx context.Context, client alago.AlagoClient, in *EventNextInput) (*EventNextOutput, error) {
	if in == nil {
		return nil, fmt.Errorf("EventNextInput is nil")
	}
	if in.LambdaExtensionIdentifier == "" {
		return nil, fmt.Errorf("EventNextInput.LambdaExtensionIdentifier is empty")
	}
	hs := []internal.Header{
		{Key: requestHeaderLambdaExtensionIdentifier, Value: in.LambdaExtensionIdentifier},
	}

	url := fmt.Sprintf(eventNextEndpointFmt, client.Host())
	sc, h, b, err := internal.CallAPI(context.Background(), client, http.MethodGet, url, nil, hs...)
	if err != nil {
		return nil, err
	}

	out, err := generateEventNextOutput(sc, h, b)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func generateEventNextOutput(sc int, header http.Header, body []byte) (*EventNextOutput, error) {
	out := EventNextOutput{}
	out.StatusCode = sc

	if v, ok := header[responseHeaderLambdaExtensionEventIdentifier]; ok && len(v) > 0 {
		out.LambdaExtensionEventIdentifier = v[0]
	}

	if sc != http.StatusOK {
		var errRes ErrorResponse
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}
		out.Error = &errRes
		return &out, nil
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("err:%v, body:%s", err, string(body))
	}

	return &out, nil
}
