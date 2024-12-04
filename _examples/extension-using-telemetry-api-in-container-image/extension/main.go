package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"telemetry-api-extension-exemple/extension"
	"telemetry-api-extension-exemple/logger"
	"telemetry-api-extension-exemple/telemetry"
)

var extensionName = filepath.Base(os.Args[0])

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	l := logger.NewLogger(os.Stdout, "Telemetry API Extension Main")

	go func() {
		s := <-sigs
		cancel()

		l.Info("Received Signal: %v", s)
		l.Info("Exiting")
	}()

	hce := &http.Client{Timeout: 0}
	le := logger.NewLogger(os.Stdout, "Telemetry API Extension Client")
	extensionClient, err := extension.NewClient(hce, le)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	// register extension
	if err := extensionClient.Register(ctx, extensionName); err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	// start Telemetry API Subscriber
	lts := logger.NewLogger(os.Stdout, "Telemetry API Subscriber")
	subscriber := telemetry.NewTelemetryAPISubscriber(lts)
	address, err := subscriber.Start()
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	// subscribe Telemetry API
	ltc := logger.NewLogger(os.Stdout, "Telemetry API Client")
	telemetryClient, err := telemetry.NewClient(http.DefaultClient, ltc)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	if err := telemetryClient.Subscribe(ctx, extensionClient.ExtensionIdentifier(), address); err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	if err := processEvents(ctx, extensionClient); err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
}

func processEvents(ctx context.Context, c *extension.Client) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if waitNextEvent, err := c.PollingEvent(ctx); err != nil {
				return err
			} else if !waitNextEvent {
				// received shutdown event
				return nil
			}
		}
	}
}
