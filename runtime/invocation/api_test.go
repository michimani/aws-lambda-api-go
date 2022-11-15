package invocation_test

import (
	"net/http"
	"testing"

	"github.com/michimani/aws-lambda-api-go/runtime"
	"github.com/michimani/aws-lambda-api-go/runtime/invocation"
	"github.com/stretchr/testify/assert"
)

func Test_generateNextOutput(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		header     http.Header
		body       []byte
		expect     *invocation.NextOutput
		wantErr    bool
	}{
		{
			name:       "ok",
			statusCode: 200,
			header: map[string][]string{
				"Lambda-Runtime-Aws-Request-Id":       {"lambda-runtime-aws-request-id"},
				"Lambda-Runtime-Trace-Id":             {"lambda-runtime-trace-id"},
				"Lambda-Runtime-Client-Context":       {"lambda-runtime-client-context"},
				"Lambda-Runtime-Cognito-Identity":     {"lambda-runtime-cognito-identity"},
				"Lambda-Runtime-Deadline-Ms":          {"lambda-runtime-deadline-ms"},
				"Lambda-Runtime-Invoked-Function-Arn": {"lambda-runtime-invoked-function-arn"},
			},
			body: []byte("test-response-body"),
			expect: &invocation.NextOutput{
				StatusCode:         200,
				AWSRequestID:       "lambda-runtime-aws-request-id",
				TraceID:            "lambda-runtime-trace-id",
				ClientContext:      "lambda-runtime-client-context",
				CognitoIdentity:    "lambda-runtime-cognito-identity",
				DeadlineMs:         "lambda-runtime-deadline-ms",
				InvokedFunctionArn: "lambda-runtime-invoked-function-arn",
				RawEventResponse:   []byte("test-response-body"),
			},
			wantErr: false,
		},
		{
			name:       "ok: not OK status code",
			statusCode: 403,
			header: map[string][]string{
				"Lambda-Runtime-Aws-Request-Id":       {"lambda-runtime-aws-request-id"},
				"Lambda-Runtime-Trace-Id":             {"lambda-runtime-trace-id"},
				"Lambda-Runtime-Client-Context":       {"lambda-runtime-client-context"},
				"Lambda-Runtime-Cognito-Identity":     {"lambda-runtime-cognito-identity"},
				"Lambda-Runtime-Deadline-Ms":          {"lambda-runtime-deadline-ms"},
				"Lambda-Runtime-Invoked-Function-Arn": {"lambda-runtime-invoked-function-arn"},
			},
			body: []byte(`{"errorMessage":"test-error-message", "errorType":"test-error-type"}`),
			expect: &invocation.NextOutput{
				StatusCode: 403,
				Error: &runtime.ErrorResponse{
					ErrorMessage: "test-error-message",
					ErrorType:    "test-error-type",
				},
			},
			wantErr: false,
		},
		{
			name:       "ng: failed to unmarshal error response",
			statusCode: 403,
			header: map[string][]string{
				"Lambda-Runtime-Aws-Request-Id":       {"lambda-runtime-aws-request-id"},
				"Lambda-Runtime-Trace-Id":             {"lambda-runtime-trace-id"},
				"Lambda-Runtime-Client-Context":       {"lambda-runtime-client-context"},
				"Lambda-Runtime-Cognito-Identity":     {"lambda-runtime-cognito-identity"},
				"Lambda-Runtime-Deadline-Ms":          {"lambda-runtime-deadline-ms"},
				"Lambda-Runtime-Invoked-Function-Arn": {"lambda-runtime-invoked-function-arn"},
			},
			body:    []byte(`///`),
			expect:  nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			out, err := invocation.Exported_generateNextOutput(c.statusCode, c.header, c.body)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(out)
				return
			}

			asst.NoError(err)
			asst.Equal(*c.expect, *out)
		})
	}
}
