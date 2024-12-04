package extension

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"telemetry-api-extension-exemple/logger"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/extension"
)

type Client struct {
	alagoClient         *alago.Client
	logger              *logger.Logger
	extensionIdentifier string
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

func (c *Client) ExtensionIdentifier() string {
	return c.extensionIdentifier
}

func (c *Client) Register(ctx context.Context, extensionName string) error {
	out, err := extension.Register(ctx, c.alagoClient, &extension.RegisterInput{
		LambdaExtensionName: extensionName,
		Events: []extension.EventType{
			extension.EventTypeInvoke,
			extension.EventTypeShutdown,
		},
	})

	if err != nil {
		return fmt.Errorf("An expected error occurred at extension registration. err:%v", err)
	}

	if out.StatusCode != http.StatusOK {
		return fmt.Errorf("An error occurred at extension registration. statusCode:%d errType:%s errMessage:%s",
			out.StatusCode, out.Error.ErrorType, out.Error.ErrorMessage)
	}

	c.logger.Info("Succeeded to register extension.")
	c.extensionIdentifier = out.LambdaExtensionIdentifier

	return nil
}

const (
	eventTypeInvoke   string = "INVOKE"
	eventTypeShutdown string = "SHUTDOWN"
)

// PollingEvent call GET /extension/event/next and handle event.
// returns
// - bool: if continue to polling, returns true
// - error
func (c *Client) PollingEvent(ctx context.Context) (bool, error) {
	c.logger.Info("Waiting for next event...")
	out, err := extension.EventNext(ctx, c.alagoClient, &extension.EventNextInput{
		LambdaExtensionIdentifier: c.extensionIdentifier,
	})

	if err != nil {
		return false, err
	}

	if out.StatusCode != http.StatusOK {
		return false, fmt.Errorf("An error occurred at calling /extension/event/next API. statusCode:%d errType:%s errMessage:%s",
			out.StatusCode, out.Error.ErrorType, out.Error.ErrorMessage)
	}

	switch out.EventType {
	case eventTypeInvoke:
		now := time.Now().UTC()
		c.logger.Info("Received invoke event. awsRequestId:%s invokedAt:%v", out.RequestID, now)
	case eventTypeShutdown:
		c.logger.Info("Received shutdown event. reason:%s", out.ShutdownReason)
		return false, nil
	default:
		return false, fmt.Errorf("Cannot handle event. eventType:%s", out.EventType)
	}

	return true, nil
}
