package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/internal"
)

const (
	subscribeEndpointFmt string = "http://%s/2022-07-01/telemetry"
)

// To subscribe to a telemetry stream, a Lambda extension can send a Subscribe API request.
//
// https://docs.aws.amazon.com/lambda/latest/dg/telemetry-api-reference.html#telemetry-subscribe-api
func Subscribe(ctx context.Context, client alago.AlagoClient, in *SubscribeInput) (*SubscribeOutput, error) {
	if in == nil {
		return nil, fmt.Errorf("SubscribeInput is nil")
	}

	reqBody, err := inputToRequestBody(in)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(subscribeEndpointFmt, client.Host())
	sc, h, b, err := internal.CallAPI(context.Background(), client, http.MethodPut, url, reqBody)
	if err != nil {
		return nil, err
	}

	out, err := generateSubscribeOutput(sc, h, b)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func generateSubscribeOutput(sc int, header http.Header, body []byte) (*SubscribeOutput, error) {
	out := SubscribeOutput{}
	out.StatusCode = sc

	if sc != http.StatusOK {
		var errRes ErrorResponse
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}
		out.Error = &errRes
		return &out, nil
	}

	return &out, nil
}
