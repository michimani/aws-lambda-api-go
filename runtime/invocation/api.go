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
