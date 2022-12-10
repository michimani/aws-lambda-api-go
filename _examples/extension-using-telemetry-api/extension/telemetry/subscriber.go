package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"telemetry-api-extension-exemple/logger"
	"time"
)

const defaultSubscriberPort = "1210"
const address = "sandbox.localdomain:" + defaultSubscriberPort

type TelemetryAPISubscriber struct {
	httpServer *http.Server
	logger     *logger.Logger
}

func NewTelemetryAPISubscriber(l *logger.Logger) *TelemetryAPISubscriber {
	return &TelemetryAPISubscriber{
		httpServer: nil,
		logger:     l,
	}
}

func (s *TelemetryAPISubscriber) Start() (string, error) {
	s.logger.Info("Starting on address:%s", address)
	s.httpServer = &http.Server{Addr: address}
	http.HandleFunc("/", s.telemetryEventHandler)

	go func() {
		err := s.httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			s.logger.Error("Unexpected stop on Http Server. err:%v", err)
			s.Shutdown()
		} else {
			s.logger.Info("Http Server closed. err:%v", err)
		}
	}()

	return fmt.Sprintf("http://%s/", address), nil
}

type telemetryAPIEvent struct {
	Record any       `json:"record"`
	Type   string    `json:"type"`
	Time   time.Time `json:"time"`
}

func (s *TelemetryAPISubscriber) telemetryEventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("Failed to reading body. err:%v", err)
		return
	}

	var events []telemetryAPIEvent
	_ = json.Unmarshal(body, &events)

	s.logger.Info("Received %d events.", len(events))
	for i, e := range events {
		s.logger.Info("%d: Time:%s Type:%s Record:%v", i, e.Time.Format(time.RFC3339Nano), e.Type, e.Record)
	}

	events = nil
}

func (s *TelemetryAPISubscriber) Shutdown() {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			s.logger.Error("Failed to shutdown http server gracefully:%v", err)
		} else {
			s.httpServer = nil
		}
	}
}
