package telemetry_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/michimani/aws-lambda-api-go/telemetry"
	"github.com/stretchr/testify/assert"
)

func Test_DestinationProtocol_Valid(t *testing.T) {
	cases := []struct {
		name   string
		dp     telemetry.DestinationProtocol
		expect bool
	}{
		{
			name:   "HTTP",
			dp:     telemetry.DestinationProtocolHTTP,
			expect: true,
		},
		{
			name:   "TCP",
			dp:     telemetry.DestinationProtocolTCP,
			expect: true,
		},
		{
			name:   "invalid value",
			dp:     telemetry.DestinationProtocol("invalid value"),
			expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			asst.Equal(c.expect, c.dp.Valid())
		})
	}
}

func Test_TelemetryType_Valid(t *testing.T) {
	cases := []struct {
		name   string
		tlt    telemetry.TelemetryType
		expect bool
	}{
		{
			name:   "platform",
			tlt:    telemetry.TelemetryTypePlatform,
			expect: true,
		},
		{
			name:   "function",
			tlt:    telemetry.TelemetryTypeFunction,
			expect: true,
		},
		{
			name:   "extension",
			tlt:    telemetry.TelemetryTypeExtension,
			expect: true,
		},
		{
			name:   "invalid value",
			tlt:    telemetry.TelemetryType("invalid value"),
			expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			asst.Equal(c.expect, c.tlt.Valid())
		})
	}
}

func Test_inputToRequestBody(t *testing.T) {
	var (
		bufMaxItems     = 100
		bufMaxBytes     = 512 * 1024
		bufMaxTimeoutMs = 100
	)

	cases := []struct {
		name    string
		in      *telemetry.SubscribeInput
		expect  string
		wantErr bool
	}{
		{
			name: "ok: with default buffering",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocolHTTP,
				DestinationURI:      "localhost",
				TelemetryTypes:      []telemetry.TelemetryType{telemetry.TelemetryTypePlatform},
			},
			expect:  `{"schemaVersion":"2022-07-01","destination":{"protocol":"HTTP","URI":"localhost"},"types":["platform"],"buffering":{"maxItems":1000,"maxBytes":262144,"timeoutMs":10000}}`,
			wantErr: false,
		},
		{
			name: "ok: with custom buffering",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocolHTTP,
				DestinationURI:      "localhost",
				TelemetryTypes:      []telemetry.TelemetryType{telemetry.TelemetryTypePlatform},
				BufferMaxItems:      &bufMaxItems,
				BufferMaxBytes:      &bufMaxBytes,
				BufferTimeoutMs:     &bufMaxTimeoutMs,
			},
			expect:  `{"schemaVersion":"2022-07-01","destination":{"protocol":"HTTP","URI":"localhost"},"types":["platform"],"buffering":{"maxItems":100,"maxBytes":524288,"timeoutMs":100}}`,
			wantErr: false,
		},
		{
			name: "ng: invalid DestinationProtocol value",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocol("invalid value"),
				DestinationURI:      "localhost",
				TelemetryTypes:      []telemetry.TelemetryType{telemetry.TelemetryTypePlatform},
			},
			expect:  ``,
			wantErr: true,
		},
		{
			name: "ng: DestinationURI is empty",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocolHTTP,
				DestinationURI:      "",
				TelemetryTypes:      []telemetry.TelemetryType{telemetry.TelemetryTypePlatform},
			},
			expect:  ``,
			wantErr: true,
		},
		{
			name: "ng: TelemetryTypes is empty",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocolHTTP,
				DestinationURI:      "localhost",
				TelemetryTypes:      []telemetry.TelemetryType{},
			},
			expect:  ``,
			wantErr: true,
		},
		{
			name: "ng: TelemetryTypes includes invalid value",
			in: &telemetry.SubscribeInput{
				DestinationProtocol: telemetry.DestinationProtocolHTTP,
				DestinationURI:      "localhost",
				TelemetryTypes: []telemetry.TelemetryType{
					telemetry.TelemetryTypeExtension,
					telemetry.TelemetryType("invalid vale"),
				},
			},
			expect:  ``,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			body, err := telemetry.Exported_inputToRequestBody(c.in)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(body)
				return
			}

			asst.NoError(err)
			asst.NotNil(body)

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, body)
			asst.NoError(err)

			asst.Equal(c.expect, buf.String())
		})
	}
}
