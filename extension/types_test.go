package extension_test

import (
	"io"
	"testing"

	"github.com/michimani/aws-lambda-api-go/extension"
	"github.com/stretchr/testify/assert"
)

func Test_event_toRequestBody_error(t *testing.T) {
	cases := []struct {
		name    string
		es      *extension.Exported_event
		expect  io.Reader
		wantErr bool
	}{
		{
			name:    "error: es is nil",
			es:      nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			r, err := c.es.Exported_toRequestBody()
			if c.wantErr {
				asst.Error(err)
				asst.Nil(r)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, r)
		})
	}
}
