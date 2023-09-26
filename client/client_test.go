package client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/oasdiff/telemetry/client"
	"github.com/oasdiff/telemetry/model"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type DiffFlags struct {
	base                     string
	revision                 string
	composed                 bool
	prefixBase               string
	prefixRevision           string
	stripPrefixBase          string
	stripPrefixRevision      string
	matchPath                string
	filterExtension          string
	failOnDiff               bool
	circularReferenceCounter int
	includePathParams        bool
}

func TestSend(t *testing.T) {

	const subCommand, version = "diff", "v1.2.3"

	cmd := &cobra.Command{
		Use:   model.Application,
		Short: "OpenAPI specification diff",
	}
	cmd.Version = version
	cmd.SetArgs([]string{subCommand, "https://aerial-data-production.herokuapp.com/bank/api/openapi3.json", "https://app.swaggerhub.com/apis/g4/Banking/1.7.56", "--composed", "--max-circular-dep", "7", "--match-path", "a/b/c"})
	cmd.AddCommand(getDiffCmd())
	require.NoError(t, cmd.Execute())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)

		var events map[string][]*model.Telemetry
		require.NoError(t, json.NewDecoder(r.Body).Decode(&events))
		telemetry := events[model.KeyEvents][0]
		require.True(t, time.Now().UnixMilli()-telemetry.Time.UnixMilli() < 100000)
		require.NotEmpty(t, telemetry.MachineId)
		require.NotEmpty(t, telemetry.Runtime)
		require.NotEmpty(t, telemetry.Platform)
		require.Equal(t, subCommand, telemetry.Command)
		require.Equal(t, version, telemetry.ApplicationVersion)
		require.Len(t, telemetry.Args, 2)
		require.Equal(t, "heroku", telemetry.Args[0])
		require.Equal(t, "swaggerhub", telemetry.Args[1])
		require.Len(t, telemetry.Flags, 3)
		require.Equal(t, "true", telemetry.Flags["composed"])
		require.Equal(t, "[redacted]", telemetry.Flags["match-path"])
		require.Equal(t, "7", telemetry.Flags["max-circular-dep"])
	}))

	c := client.NewDefaultCollector()
	c.EventsUrl = server.URL
	c.SendCommand(cmd)
}

func getDiffCmd() *cobra.Command {

	flags := DiffFlags{}

	cmd := cobra.Command{
		Use:   "diff base revision [flags]",
		Short: "Generate a diff report",
		Long: `Generate a diff report between base and revision specs.
Base and revision can be a path to a file or a URL.
In 'composed' mode, base and revision can be a glob and oasdiff will compare mathcing endpoints between the two sets of files.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.base = args[0]
			flags.revision = args[1]
			cmd.Root().SilenceUsage = true
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&flags.composed, "composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	cmd.PersistentFlags().StringVarP(&flags.matchPath, "match-path", "p", "", "include only paths that match this regular expression")
	cmd.PersistentFlags().StringVarP(&flags.filterExtension, "filter-extension", "", "", "exclude paths and operations with an OpenAPI Extension matching this regular expression")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	cmd.PersistentFlags().StringVarP(&flags.prefixBase, "prefix-base", "", "", "add this prefix to paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.prefixRevision, "prefix-revision", "", "", "add this prefix to paths in revised-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixBase, "strip-prefix-base", "", "", "strip this prefix from paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixRevision, "strip-prefix-revision", "", "", "strip this prefix from paths in revised-spec before comparison")
	cmd.PersistentFlags().BoolVarP(&flags.includePathParams, "include-path-params", "", false, "include path parameter names in endpoint matching")
	cmd.PersistentFlags().BoolVarP(&flags.failOnDiff, "fail-on-diff", "o", false, "exit with return code 1 when any change is found")

	return &cmd
}
