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
	eventNextEndpointFmt string = "http://%s/2020-01-01/extension/event/next"

	responseHeaderLambdaExtensionEventIdentifier string = "Lambda-Extension-Event-Identifier"
)

// Extension makes this HTTP request when it is ready to receive and process a new event.
// This is an iterator-style blocking API call. Response contains event JSON document.
//
// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-extensions-api.html#extensions-api-next
func EventNext(ctx context.Context, client alago.AlagoClient) (*EventNextOutput, error) {
	url := fmt.Sprintf(eventNextEndpointFmt, client.Host())
	sc, h, b, err := internal.CallAPI(context.Background(), client, http.MethodGet, url, nil)
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
