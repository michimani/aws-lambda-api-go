package telemetry_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/telemetry"
	"github.com/michimani/http-client-mock/hcmock"
	"github.com/stretchr/testify/assert"
)

func Test_Subscribe(t *testing.T) {
	cases := []struct {
		name       string
		httpClient *http.Client
		host       string
		in         *telemetry.SubscribeInput
		expect     *telemetry.SubscribeOutput
		wantErr    bool
	}{
		{
			name: "ok",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				BodyBytes:  []byte(`OK`),
			}),
			host: "test-host",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocolHTTP,
				DestinationURI:      "http://localhost",
				TelemetryTypes:      []telemetry.TelemetryType{telemetry.TelemetryTypePlatform},
			},
			expect: &telemetry.SubscribeOutput{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "ng: SubscribeInput is nil",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				BodyBytes:  []byte(`OK`),
			}),
			host:    "test-host",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: inputToRequestBody returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				BodyBytes:  []byte(`OK`),
			}),
			host:    "test-host",
			in:      &telemetry.SubscribeInput{},
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: CallAPI returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				BodyBytes:  []byte(`OK`),
			}),
			host: "\U00000001",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocolHTTP,
				DestinationURI:      "http://localhost",
				TelemetryTypes:      []telemetry.TelemetryType{telemetry.TelemetryTypePlatform},
			},
			expect:  nil,
			wantErr: true,
		},
		{
			name: "ng: generateSubscribeOutput returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 400,
				BodyBytes:  []byte(`///`),
			}),
			host: "test-host",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocolHTTP,
				DestinationURI:      "http://localhost",
				TelemetryTypes:      []telemetry.TelemetryType{telemetry.TelemetryTypePlatform},
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

			out, err := telemetry.Subscribe(context.Background(), ac, c.in)
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

func Test_generateEventSubscribeOutput(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		header     http.Header
		body       []byte
		expect     *telemetry.SubscribeOutput
		wantErr    bool
	}{
		{
			name:       "ok",
			statusCode: 200,
			body:       []byte(`OK`),
			expect: &telemetry.SubscribeOutput{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name:       "ok: not OK status code",
			statusCode: 400,
			body:       []byte(`{"errorMessage":"test-error-message", "errorType":"test-error-type"}`),
			expect: &telemetry.SubscribeOutput{
				StatusCode: 400,
				Error: &telemetry.ErrorResponse{
					ErrorMessage: "test-error-message",
					ErrorType:    "test-error-type",
				},
			},
			wantErr: false,
		},
		{
			name:       "ng: failed to unmarshal error response",
			statusCode: 400,
			body:       []byte(`///`),
			expect:     nil,
			wantErr:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			out, err := telemetry.Exported_generateSubscribeOutput(c.statusCode, c.header, c.body)
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
