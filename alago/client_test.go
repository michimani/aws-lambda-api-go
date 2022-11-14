package alago_test

import (
	"net/http"
	"testing"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	cases := []struct {
		name    string
		in      *alago.NewClientInput
		envKey  string
		wantErr bool
	}{
		{
			name:    "ok: use default http client",
			in:      &alago.NewClientInput{},
			envKey:  "AWS_LAMBDA_RUNTIME_API",
			wantErr: false,
		},
		{
			name: "ok: use custom http client",
			in: &alago.NewClientInput{
				HttpClient: &http.Client{
					Timeout: 300,
				},
			},
			envKey:  "AWS_LAMBDA_RUNTIME_API",
			wantErr: false,
		},
		{
			name:    "ng: NewClientInput is nil",
			in:      nil,
			envKey:  "AWS_LAMBDA_RUNTIME_API",
			wantErr: true,
		},
		{
			name:    "ng: env key is not set",
			in:      &alago.NewClientInput{},
			envKey:  "",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			if c.envKey != "" {
				tt.Setenv(c.envKey, "test-env-value")
			}

			client, err := alago.NewClient(c.in)

			if c.wantErr {
				asst.Error(err)
				asst.Nil(client)
				return
			}

			asst.NoError(err)
			asst.NotNil(client)
		})
	}
}
