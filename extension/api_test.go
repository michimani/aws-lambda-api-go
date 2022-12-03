package extension_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/extension"
	"github.com/michimani/http-client-mock/hcmock"
	"github.com/stretchr/testify/assert"
)

func Test_EventNext(t *testing.T) {
	cases := []struct {
		name       string
		httpClient *http.Client
		host       string
		in         *extension.EventNextInput
		expect     *extension.EventNextOutput
		wantErr    bool
	}{
		{
			name: "ok: event type INVOKE",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"INVOKE","deadlineMs":123456,"requestId":"aws-request-id","invokedFunctionArn":"function-arn","tracing":{"type":"X-Amzn-Trace-Id","value":"tracing-value"}}`),
			}),
			host: "test-host",
			in: &extension.EventNextInput{
				LambdaExtensionIdentifier: "test",
			},
			expect: &extension.EventNextOutput{
				StatusCode:                     200,
				LambdaExtensionEventIdentifier: "lambda-extension-event-identifier",
				EventType:                      "INVOKE",
				DeadlineMs:                     123456,
				RequestID:                      "aws-request-id",
				InvokedFunctionArn:             "function-arn",
				Tracing:                        extension.XRayTracingInfo{Type: "X-Amzn-Trace-Id", Value: "tracing-value"},
			},
			wantErr: false,
		},
		{
			name: "ok: event type SHUTDOWN",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"SHUTDOWN","shutdownReason":"test down","deadlineMs":123456}`),
			}),
			host: "test-host",
			in: &extension.EventNextInput{
				LambdaExtensionIdentifier: "test",
			},
			expect: &extension.EventNextOutput{
				StatusCode:                     200,
				LambdaExtensionEventIdentifier: "lambda-extension-event-identifier",
				EventType:                      "SHUTDOWN",
				DeadlineMs:                     123456,
				ShutdownReason:                 "test down",
			},
			wantErr: false,
		},
		{
			name: "ok: event type SHUTDOWN",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"SHUTDOWN","shutdownReason":"test down","deadlineMs":123456}`),
			}),
			host: "test-host",
			in: &extension.EventNextInput{
				LambdaExtensionIdentifier: "test",
			},
			expect: &extension.EventNextOutput{
				StatusCode:                     200,
				LambdaExtensionEventIdentifier: "lambda-extension-event-identifier",
				EventType:                      "SHUTDOWN",
				DeadlineMs:                     123456,
				ShutdownReason:                 "test down",
			},
			wantErr: false,
		},
		{
			name: "ng: EventNextInput is nil",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"SHUTDOWN","shutdownReason":"test down","deadlineMs":123456}`),
			}),
			host:    "test-host",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: EventNextInput is empty",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"SHUTDOWN","shutdownReason":"test down","deadlineMs":123456}`),
			}),
			host:    "test-host",
			in:      &extension.EventNextInput{},
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: CallAPI returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"INVOKE","deadlineMs":123456,"requestId":"aws-request-id","invokedFunctionArn":"function-arn","tracing":{"type":"X-Amzn-Trace-Id","value":"tracing-value"}}`),
			}),
			host: "\U00000001",
			in: &extension.EventNextInput{
				LambdaExtensionIdentifier: "test",
			},
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: generateNextOutput returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`///`),
			}),
			host: "test-host",
			in: &extension.EventNextInput{
				LambdaExtensionIdentifier: "test",
			},
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

			out, err := extension.EventNext(context.Background(), ac, c.in)
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

func Test_generateEventEventNextOutput(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		header     http.Header
		body       []byte
		expect     *extension.EventNextOutput
		wantErr    bool
	}{
		{
			name:       "ok: event type INVOKE",
			statusCode: 200,
			header: map[string][]string{
				"Lambda-Extension-Event-Identifier": {"lambda-extension-event-identifier"},
			},
			body: []byte(`{"eventType":"INVOKE","deadlineMs":123456,"requestId":"aws-request-id","invokedFunctionArn":"function-arn","tracing":{"type":"X-Amzn-Trace-Id","value":"tracing-value"}}`),
			expect: &extension.EventNextOutput{
				StatusCode:                     200,
				LambdaExtensionEventIdentifier: "lambda-extension-event-identifier",
				EventType:                      "INVOKE",
				DeadlineMs:                     123456,
				RequestID:                      "aws-request-id",
				InvokedFunctionArn:             "function-arn",
				Tracing:                        extension.XRayTracingInfo{Type: "X-Amzn-Trace-Id", Value: "tracing-value"},
			},
			wantErr: false,
		},
		{
			name:       "ok: event type SHUTDOWN",
			statusCode: 200,
			header: map[string][]string{
				"Lambda-Extension-Event-Identifier": {"lambda-extension-event-identifier"},
			},
			body: []byte(`{"eventType":"SHUTDOWN","shutdownReason":"test down","deadlineMs":123456}`),
			expect: &extension.EventNextOutput{
				StatusCode:                     200,
				LambdaExtensionEventIdentifier: "lambda-extension-event-identifier",
				EventType:                      "SHUTDOWN",
				DeadlineMs:                     123456,
				ShutdownReason:                 "test down",
			},
			wantErr: false,
		},
		{
			name:       "ok: not OK status code",
			statusCode: 403,
			header: map[string][]string{
				"Lambda-Extension-Event-Identifier": {"lambda-extension-event-identifier"},
			},
			body: []byte(`{"errorMessage":"test-error-message", "errorType":"test-error-type"}`),
			expect: &extension.EventNextOutput{
				StatusCode:                     403,
				LambdaExtensionEventIdentifier: "lambda-extension-event-identifier",
				Error: &extension.ErrorResponse{
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
				"Lambda-Extension-Event-Identifier": {"lambda-extension-event-identifier"},
			},
			body:    []byte(`///`),
			expect:  nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			out, err := extension.Exported_generateEventNextOutput(c.statusCode, c.header, c.body)
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
