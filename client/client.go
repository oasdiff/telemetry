package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/oasdiff/go-common/util"
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
	EventsUrl   string
	ignoreFlags *util.StringSet
}

func NewCollector(ignoreFlags *util.StringSet) *Collector {

	return &Collector{
		EventsUrl:   home.JoinPath(model.KeyEvents).String(),
		ignoreFlags: getStringSet(ignoreFlags),
	}
}

func (c *Collector) Send(cmd *cobra.Command) error {

	return send(c.EventsUrl, redact(fromCommand(cmd), c.ignoreFlags))
}

func redact(telemetry *model.Telemetry, ignoreFlags *util.StringSet) *model.Telemetry {

	return redactFlags(redactArgs(telemetry), ignoreFlags)
}

func redactFlags(telemetry *model.Telemetry, ignoreFlags *util.StringSet) *model.Telemetry {

	for key := range telemetry.Flags {
		if ignoreFlags.Has(key) {
			telemetry.Flags[key] = ""
		}
	}

	return telemetry
}

func redactArgs(telemetry *model.Telemetry) *model.Telemetry {

	for i, arg := range telemetry.Args {
		if arg != "" {
			telemetry.Args[i] = getCategory(arg)
		}
	}

	return telemetry
}

func getCategory(name string) string {

	if web := getWebCategory(name); web != "" {
		return web
	}
	return "file"
}

func getWebCategory(name string) string {

	if isSwaggerHub(name) {
		return "swaggerhub"
	} else if isGitHub(name) {
		return "github"
	} else if isGCS(name) {
		return "gcs"
	} else if isS3(name) {
		return "s3"
	} else if isAzure(name) {
		return "azure"
	} else if isHeroku(name) {
		return "heroku"
	} else if strings.HasPrefix(name, "https://") {
		return "https"
	} else if strings.HasPrefix(name, "http://") {
		return "http"
	}

	return ""
}

func isHeroku(name string) bool {

	res, err := regexp.MatchString(`^https://(.)*herokuapp\.com/`, name)
	if err != nil {
		slog.Debug("failed to validate if name is heroku host", err)
		return false
	}

	return res
}

func isGCS(name string) bool {

	res, err := regexp.MatchString(`^https://storage\.cloud\.google\.com/`, name)
	if err != nil {
		slog.Debug("failed to validate if name is GCS host", err)
		return false
	}

	return res
}

func isS3(name string) bool {

	// TODO
	return false
}

func isAzure(name string) bool {

	res, err := regexp.MatchString(`^https://(.)*azure.com/`, name)
	if err != nil {
		slog.Debug("failed to validate if name is azure host", err)
		return false
	}

	return res
}

func isGitHub(name string) bool {

	res, err := regexp.MatchString(`^https://(.)*githubusercontent.com/`, name)
	if err != nil {
		slog.Debug("failed to validate if name is githubusercontent host", err)
		return false
	}
	if !res {
		res, err = regexp.MatchString(`^https://(.)*github.com/`, name)
		if err != nil {
			slog.Debug("failed to validate if name is github host", err)
			return false
		}
	}

	return res
}

func isSwaggerHub(name string) bool {

	res, err := regexp.MatchString(`^https://(.)*swaggerhub.com/`, name)
	if err != nil {
		slog.Debug("failed to validate if name is github host", err)
		return false
	}

	return res
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

func getStringSet(ignoreFlags *util.StringSet) *util.StringSet {

	if ignoreFlags == nil {
		return util.NewStringSet()
	}
	return ignoreFlags
}
