package runtime_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/runtime"
	"github.com/michimani/http-client-mock/hcmock"
	"github.com/stretchr/testify/assert"
)

func Test_InvocationNext(t *testing.T) {
	cases := []struct {
		name       string
		httpClient *http.Client
		host       string
		expect     *runtime.NextOutput
		wantErr    bool
	}{
		{
			name: "ok",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Runtime-Aws-Request-Id", Value: "lambda-runtime-aws-request-id"},
					{Key: "Lambda-Runtime-Trace-Id", Value: "lambda-runtime-trace-id"},
					{Key: "Lambda-Runtime-Client-Context", Value: "lambda-runtime-client-context"},
					{Key: "Lambda-Runtime-Cognito-Identity", Value: "lambda-runtime-cognito-identity"},
					{Key: "Lambda-Runtime-Deadline-Ms", Value: "lambda-runtime-deadline-ms"},
					{Key: "Lambda-Runtime-Invoked-Function-Arn", Value: "lambda-runtime-invoked-function-arn"},
				},
				BodyBytes: []byte(`test-response-body`),
			}),
			host: "test-host",
			expect: &runtime.NextOutput{
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
			name: "ng: CallAPI returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Runtime-Aws-Request-Id", Value: "lambda-runtime-aws-request-id"},
					{Key: "Lambda-Runtime-Trace-Id", Value: "lambda-runtime-trace-id"},
					{Key: "Lambda-Runtime-Client-Context", Value: "lambda-runtime-client-context"},
					{Key: "Lambda-Runtime-Cognito-Identity", Value: "lambda-runtime-cognito-identity"},
					{Key: "Lambda-Runtime-Deadline-Ms", Value: "lambda-runtime-deadline-ms"},
					{Key: "Lambda-Runtime-Invoked-Function-Arn", Value: "lambda-runtime-invoked-function-arn"},
				},
				BodyBytes: []byte(`test-response-body`),
			}),
			host:    "\U00000001",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: generateNextOutput returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 403,
				Headers: []hcmock.Header{
					{Key: "Lambda-Runtime-Aws-Request-Id", Value: "lambda-runtime-aws-request-id"},
					{Key: "Lambda-Runtime-Trace-Id", Value: "lambda-runtime-trace-id"},
					{Key: "Lambda-Runtime-Client-Context", Value: "lambda-runtime-client-context"},
					{Key: "Lambda-Runtime-Cognito-Identity", Value: "lambda-runtime-cognito-identity"},
					{Key: "Lambda-Runtime-Deadline-Ms", Value: "lambda-runtime-deadline-ms"},
					{Key: "Lambda-Runtime-Invoked-Function-Arn", Value: "lambda-runtime-invoked-function-arn"},
				},
				BodyBytes: []byte(`///`),
			}),
			host:    "test-host",
			expect:  nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			tt.Setenv("AWS_LAMBDA_RUNTIME_API", c.host)

			ac, err := alago.NewClient(&alago.NewClientInput{
				HttpClient: c.httpClient,
			})

			asst.NoError(err)

			out, err := runtime.InvocationNext(context.Background(), ac)
			if c.wantErr {
				asst.Error(err, err)
				asst.Nil(out)
				return
			}

			asst.NoError(err)
			asst.NotNil(out)
			asst.Equal(*c.expect, *out)
			asst.Equal(*c.expect, *out)
			asst.Equal(*c.expect, *out)
		})
	}
}

func Test_generateNextOutput(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		header     http.Header
		body       []byte
		expect     *runtime.NextOutput
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
			expect: &runtime.NextOutput{
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
			expect: &runtime.NextOutput{
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

			out, err := runtime.Exported_generateNextOutput(c.statusCode, c.header, c.body)
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

func Test_InvocationResponse(t *testing.T) {
	cases := []struct {
		name       string
		httpClient *http.Client
		in         *runtime.ResponseInput
		host       string
		expect     *runtime.ResponseOutput
		wantErr    bool
	}{
		{
			name: "ok",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 202,
				BodyBytes:  []byte(`{"status":"test-status"}`),
			}),
			in: &runtime.ResponseInput{
				AWSRequestID: "test-request-id",
				Response:     strings.NewReader("test-response"),
			},
			host: "test-host",
			expect: &runtime.ResponseOutput{
				StatusCode: 202,
				Status:     "test-status",
			},
			wantErr: false,
		},
		{
			name: "ng: CallAPI returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 202,
				BodyBytes:  []byte(`{"status":"test-status"}`),
			}),
			in: &runtime.ResponseInput{
				AWSRequestID: "test-request-id",
				Response:     strings.NewReader("test-response"),
			},
			host:    "\U00000001",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: AWSRequestID is empty",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 202,
				BodyBytes:  []byte(`{"status":"test-status"}`),
			}),
			in: &runtime.ResponseInput{
				Response: strings.NewReader("test-response"),
			},
			host:    "test-host",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: Response is nil",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 202,
				BodyBytes:  []byte(`{"status":"test-status"}`),
			}),
			in: &runtime.ResponseInput{
				AWSRequestID: "test-request-id",
			},
			host:    "test-host",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: ResponseInput is nil",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 202,
				BodyBytes:  []byte(`{"status":"test-status"}`),
			}),
			in:      nil,
			host:    "test-host",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: generateResponseOutput returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 403,
				BodyBytes:  []byte(`///`),
			}),
			in: &runtime.ResponseInput{
				AWSRequestID: "test-request-id",
				Response:     strings.NewReader("test-response"),
			},
			host:    "test-host",
			expect:  nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			tt.Setenv("AWS_LAMBDA_RUNTIME_API", c.host)

			ac, err := alago.NewClient(&alago.NewClientInput{
				HttpClient: c.httpClient,
			})

			asst.NoError(err)

			out, err := runtime.InvocationResponse(context.Background(), ac, c.in)
			if c.wantErr {
				asst.Error(err, err)
				asst.Nil(out)
				return
			}

			asst.NoError(err)
			asst.NotNil(out)
			asst.Equal(*c.expect, *out)
			asst.Equal(*c.expect, *out)
			asst.Equal(*c.expect, *out)
		})
	}
}

func Test_generateResponseOutput(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		body       []byte
		expect     *runtime.ResponseOutput
		wantErr    bool
	}{
		{
			name:       "ok",
			statusCode: 202,
			body:       []byte(`{"status":"test-status"}`),
			expect: &runtime.ResponseOutput{
				StatusCode: 202,
				Status:     "test-status",
			},
			wantErr: false,
		},
		{
			name:       "ok: not OK status code",
			statusCode: 403,
			body:       []byte(`{"errorMessage":"test-error-message", "errorType":"test-error-type"}`),
			expect: &runtime.ResponseOutput{
				StatusCode: 403,
				Error: &runtime.ErrorResponse{
					ErrorMessage: "test-error-message",
					ErrorType:    "test-error-type",
				},
			},
			wantErr: false,
		},
		{
			name:       "ng: failed to unmarshal ok response",
			statusCode: 202,
			body:       []byte(`///`),
			expect:     nil,
			wantErr:    true,
		},
		{
			name:       "ng: failed to unmarshal error response",
			statusCode: 403,
			body:       []byte(`///`),
			expect:     nil,
			wantErr:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			out, err := runtime.Exported_generateResponseOutput(c.statusCode, c.body)
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
