package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oasdiff/telemetry/model"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func Send(cmd *cobra.Command) error {

	return send(model.FromCommand(cmd))
}

func send(t *model.Telemetry) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(t)
	if err != nil {
		slog.Debug("failed to encode telemetry", "error", err)
		return err
	}

	response, err := http.Post("", "application/json", &buf)
	if err != nil {
		slog.Debug("failed to send telemetry", "error", err)
		return err
	}

	if response.StatusCode != http.StatusCreated {
		err := fmt.Errorf("failed to send telemetry with response status '%s'", response.Status)
		slog.Debug(err.Error())
		return err
	}

	return nil
}
