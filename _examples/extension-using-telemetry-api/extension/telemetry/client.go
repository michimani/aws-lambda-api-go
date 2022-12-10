package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"telemetry-api-extension-exemple/logger"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/telemetry"
)

type Client struct {
	alagoClient *alago.Client
	logger      *logger.Logger
}

func NewClient(hc *http.Client, l *logger.Logger) (*Client, error) {
	ac, err := alago.NewClient(&alago.NewClientInput{
		HttpClient: hc,
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		alagoClient: ac,
		logger:      l,
	}, nil
}

func (c *Client) Subscribe(ctx context.Context, exId string, httpURI string) error {
	bufTimeoutMs := 100

	out, err := telemetry.Subscribe(ctx, c.alagoClient, &telemetry.SubscribeInput{
		LambdaExtensionIdentifier: exId,
		DestinationProtocol:       telemetry.DestinationProtocolHTTP,
		DestinationURI:            httpURI,
		TelemetryTypes: []telemetry.TelemetryType{
			telemetry.TelemetryTypeFunction,
			telemetry.TelemetryTypePlatform,
		},
		BufferTimeoutMs: &bufTimeoutMs,
	})

	if err != nil {
		c.logger.Error("An error occurred at telemetrySubscribe. err:%v", err)
		return err
	}

	if out.StatusCode != http.StatusOK {
		return fmt.Errorf("An error occurred at extension registration. statusCode:%d errType:%s errMessage:%s",
			out.StatusCode, out.Error.ErrorType, out.Error.ErrorMessage)
	}

	return nil
}
