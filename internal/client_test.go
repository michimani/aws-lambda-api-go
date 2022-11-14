package internal_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/internal"
	"github.com/stretchr/testify/assert"
)

func Test_CallAPI(t *testing.T) {
	type expect struct {
		statusCode int
		header     map[string][]string
		body       []byte
	}

	cases := []struct {
		name       string
		httpClient *http.Client
		method     string
		url        string
		expect     expect
		wantErr    bool
	}{
		{
			name: "ok",
			httpClient: newMockHTTPClient(&mockInput{
				ResponseStatusCode: 200,
				ResponseHeader: map[string][]string{
					"Test-Header-Name": {"test-header-value"},
				},
				ResponseBody: io.NopCloser(strings.NewReader(`{"message": "test"}`)),
			}),
			method: "GET",
			url:    "https://example.com",
			expect: expect{
				statusCode: 200,
				header: map[string][]string{
					"Test-Header-Name": {"test-header-value"},
				},
				body: []byte(`{"message": "test"}`),
			},
			wantErr: false,
		},
		{
			name: "ok: empty response body",
			httpClient: newMockHTTPClient(&mockInput{
				ResponseStatusCode: 200,
				ResponseHeader: map[string][]string{
					"Test-Header-Name": {"test-header-value"},
				},
				ResponseBody: nil,
			}),
			method: "GET",
			url:    "https://example.com",
			expect: expect{
				statusCode: 200,
				header: map[string][]string{
					"Test-Header-Name": {"test-header-value"},
				},
				body: nil,
			},
			wantErr: false,
		},
		{
			name: "ng: failed to create request",
			httpClient: newMockHTTPClient(&mockInput{
				ResponseStatusCode: 200,
				ResponseHeader: map[string][]string{
					"Test-Header-Name": {"test-header-value"},
				},
				ResponseBody: io.NopCloser(strings.NewReader(`{"message": "test"}`)),
			}),
			method: "invalid method",
			url:    "https://example.com",
			expect: expect{
				statusCode: 200,
				header: map[string][]string{
					"Test-Header-Name": {"test-header-value"},
				},
				body: []byte(`{"message": "test"}`),
			},
			wantErr: true,
		},
		{
			name: "ng: failed to do request",
			httpClient: newMockHTTPClient(&mockInput{
				ResponseStatusCode: 200,
				ResponseHeader: map[string][]string{
					"Test-Header-Name": {"test-header-value"},
				},
				ResponseBody: io.NopCloser(strings.NewReader(`{"message": "test"}`)),
			}),
			method: "GET",
			url:    "\U00000000",
			expect: expect{
				statusCode: 200,
				header: map[string][]string{
					"Test-Header-Name": {"test-header-value"},
				},
				body: []byte(`{"message": "test"}`),
			},
			wantErr: true,
		},
	}

	t.Setenv("AWS_LAMBDA_RUNTIME_API", "test-env-value")

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			ac, err := alago.NewClient(&alago.NewClientInput{
				HttpClient: c.httpClient,
			})

			asst.NoError(err)

			sc, h, b, err := internal.CallAPI(context.Background(), ac, c.method, c.url, nil)
			if c.wantErr {
				asst.Error(err, err)
				asst.Equal(0, sc)
				asst.Nil(h)
				asst.Nil(b)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect.statusCode, sc)
			asst.Equal(c.expect.header, h)
			asst.Equal(c.expect.body, b)
		})
	}
}
