package runtime_test

import (
	"testing"

	"github.com/michimani/aws-lambda-api-go/runtime"
	"github.com/stretchr/testify/assert"
)

func Test_UnmarshalEventResponse(t *testing.T) {
	type testEvent struct {
		EventName string `json:"eventName"`
		Count     int    `json:"count"`
		IsTest    bool   `json:"isTest"`
	}

	cases := []struct {
		name    string
		o       *runtime.NextOutput
		target  *testEvent
		expect  testEvent
		wantErr bool
	}{
		{
			name: "ok",
			o: &runtime.NextOutput{
				RawEventResponse: []byte(`{"eventName":"test", "count":100, "isTest":true}`),
			},
			target: &testEvent{},
			expect: testEvent{
				EventName: "test",
				Count:     100,
				IsTest:    true,
			},
			wantErr: false,
		},
		{
			name:   "ng: receiver is nil",
			o:      nil,
			target: &testEvent{},
			expect: testEvent{
				EventName: "test",
				Count:     100,
				IsTest:    true,
			},
			wantErr: true,
		},
		{
			name: "ng: failed to unmarshal json",
			o: &runtime.NextOutput{
				RawEventResponse: []byte(`///`),
			},
			target: nil,
			expect: testEvent{
				EventName: "test",
				Count:     100,
				IsTest:    true,
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			err := c.o.UnmarshalEventResponse(c.target)
			if c.wantErr {
				asst.Error(err)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, *c.target)
		})
	}
}
