package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/oasdiff/telemetry/model"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slog"
)

var (
	home = &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("telemetry.%s.com", model.Application),
	}
)

type Collector struct {
	EventsUrl string
}

func NewCollector() *Collector {

	return &Collector{EventsUrl: home.JoinPath(model.KeyEvents).String()}
}

func (c *Collector) Send(cmd *cobra.Command) error {

	return send(c.EventsUrl, fromCommand(cmd))
}

func fromCommand(cmd *cobra.Command) *model.Telemetry {

	subCommandName := ""
	args := []string{}
	flagNameToValue := make(map[string]string)

	for _, currSubCommand := range cmd.Commands() {
		subCommandName = currSubCommand.CalledAs()
		if subCommandName != "" {
			currSubCommand.Flags().Visit(func(flag *pflag.Flag) {
				flagNameToValue[flag.Name] = flag.Value.String()
			})
			args = currSubCommand.Flags().Args()
			break
		}
	}

	return model.NewTelemetry(fmt.Sprintf("%s-cli", model.Application),
		cmd.Version, subCommandName, args, flagNameToValue)
}

func send(url string, t *model.Telemetry) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(map[string][]*model.Telemetry{"events": {t}})
	if err != nil {
		slog.Debug("failed to encode telemetry", "error", err)
		return err
	}

	response, err := http.Post(url, "application/json", &buf)
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
