package invocation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/internal"
	"github.com/michimani/aws-lambda-api-go/runtime"
)

const (
	invocationNextEndpointFmt     string = "http://%s/2018-06-01/runtime/invocation/next"
	invocationResponseEndpointFmt string = "http://%s/2018-06-01/runtime/invocation/%s/response"

	// Response Header Names
	responseHeaderLambdaRuntimeAwsRequestId       string = "Lambda-Runtime-Aws-Request-Id"
	responseHeaderLambdaRuntimeTraceId            string = "Lambda-Runtime-Trace-Id"
	responseHeaderLambdaRuntimeClientContext      string = "Lambda-Runtime-Client-Context"
	responseHeaderLambdaRuntimeCognitoIdentity    string = "Lambda-Runtime-Cognito-Identity"
	responseHeaderLambdaRuntimeDeadlineMs         string = "Lambda-Runtime-Deadline-Ms"
	responseHeaderLambdaRuntimeInvokedFunctionArn string = "Lambda-Runtime-Invoked-Function-Arn"
)

// Runtime makes this HTTP request when it is ready to receive and process a new invoke.
//
// document: https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html#runtimes-api-next
func InvocationNext(ctx context.Context, client alago.AlagoClient) (*NextOutput, error) {
	url := fmt.Sprintf(invocationNextEndpointFmt, client.Host())
	sc, h, b, err := internal.CallAPI(context.Background(), client, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	out, err := generateNextOutput(sc, h, b)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func generateNextOutput(sc int, header http.Header, body []byte) (*NextOutput, error) {
	out := NextOutput{}
	out.StatusCode = sc

	if sc != http.StatusOK {
		var errRes runtime.ErrorResponse
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}
		out.Error = &errRes
		return &out, nil
	}

	out.AWSRequestID = header.Get(responseHeaderLambdaRuntimeAwsRequestId)
	out.TraceID = header.Get(responseHeaderLambdaRuntimeTraceId)
	out.ClientContext = header.Get(responseHeaderLambdaRuntimeClientContext)
	out.CognitoIdentity = header.Get(responseHeaderLambdaRuntimeCognitoIdentity)
	out.DeadlineMs = header.Get(responseHeaderLambdaRuntimeDeadlineMs)
	out.InvokedFunctionArn = header.Get(responseHeaderLambdaRuntimeInvokedFunctionArn)
	out.RawEventResponse = body

	return &out, nil
}

// Runtime makes this request in order to submit a response.
//
// document: https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html#runtimes-api-response
func InvocationResponse(ctx context.Context, client alago.AlagoClient, in *ResponseInput) (*ResponseOutput, error) {
	if in == nil {
		return nil, fmt.Errorf("ResponseInput is nil")
	}
	if in.AWSRequestID == "" {
		return nil, fmt.Errorf("ResponseInput.AWSRequestID is empty")
	}
	if in.Response == nil {
		return nil, fmt.Errorf("ResponseInput.Response is nil")
	}

	url := fmt.Sprintf(invocationResponseEndpointFmt, client.Host(), in.AWSRequestID)
	sc, _, b, err := internal.CallAPI(context.Background(), client, http.MethodPost, url, in.Response)
	if err != nil {
		return nil, err
	}

	out, err := generateResponseOutput(sc, b)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func generateResponseOutput(sc int, body []byte) (*ResponseOutput, error) {
	out := ResponseOutput{}
	out.StatusCode = sc

	if sc != http.StatusAccepted {
		var errRes runtime.ErrorResponse
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
