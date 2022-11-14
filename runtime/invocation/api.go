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
	invocationNextEndpointFmt string = "http://%s/2018-06-01/runtime/invocation/next"

	// Response Header Names
	responseHeaderLambdaRuntimeAwsRequestId       string = "Lambda-Runtime-Aws-Request-Id"
	responseHeaderLambdaRuntimeTraceId            string = "Lambda-Runtime-Trace-Id"
	responseHeaderLambdaRuntimeClientContext      string = "Lambda-Runtime-Client-Context"
	responseHeaderLambdaRuntimeCognitoIdentity    string = "Lambda-Runtime-Cognito-Identity"
	responseHeaderLambdaRuntimeDeadlineMs         string = "Lambda-Runtime-Deadline-Ms"
	responseHeaderLambdaRuntimeInvokedFunctionArn string = "Lambda-Runtime-Invoked-Function-Arn"
)

func InvocationNext(ctx context.Context, client alago.AlagoClient) (*NextOutput, error) {
	url := fmt.Sprintf(invocationNextEndpointFmt, client.Host())
	sc, h, b, err := internal.CallAPI(context.Background(), client, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	out := NextOutput{}
	out.StatusCode = sc

	if sc != http.StatusOK {
		var errRes runtime.ErrorResponse
		if err := json.Unmarshal(b, errRes); err != nil {
			return nil, err
		}
		out.Error = &errRes
		return &out, nil
	}

	if v, ok := h[responseHeaderLambdaRuntimeAwsRequestId]; ok && len(v) > 0 {
		out.AWSRequestID = &v[0]
	}
	if v, ok := h[responseHeaderLambdaRuntimeTraceId]; ok && len(v) > 0 {
		out.TraceID = &v[0]
	}
	if v, ok := h[responseHeaderLambdaRuntimeClientContext]; ok && len(v) > 0 {
		out.ClientContext = &v[0]
	}
	if v, ok := h[responseHeaderLambdaRuntimeCognitoIdentity]; ok && len(v) > 0 {
		out.CognitoIdentity = &v[0]
	}
	if v, ok := h[responseHeaderLambdaRuntimeDeadlineMs]; ok && len(v) > 0 {
		out.DeadlineMs = &v[0]
	}
	if v, ok := h[responseHeaderLambdaRuntimeInvokedFunctionArn]; ok && len(v) > 0 {
		out.InvokedFunctionArn = &v[0]
	}

	return &out, nil
}
