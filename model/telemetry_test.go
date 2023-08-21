package model_test

import (
	"testing"

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
	format                   string
	failOnDiff               bool
	circularReferenceCounter int
	includePathParams        bool
	excludeElements          []string
}

func TestTelemetryFromCommand(t *testing.T) {

	const subCommand, version = "diff", "v1.2.3"

	cmd := &cobra.Command{
		Use:   "oasdiff",
		Short: "OpenAPI specification diff",
	}
	cmd.Version = version
	cmd.SetArgs([]string{subCommand, "data/openapi-test1.yaml", "data/openapi-test3.yaml", "--composed"})
	cmd.AddCommand(getDiffCmd())
	require.NoError(t, cmd.Execute())

	telemetry := model.FromCommand(cmd)

	require.NotEmpty(t, telemetry.Id)
	require.True(t, telemetry.Time > 0)
	require.NotEmpty(t, telemetry.MachineId)
	require.NotEmpty(t, telemetry.Runtime)
	require.NotEmpty(t, telemetry.Platform)
	require.Equal(t, subCommand, telemetry.Command)
	require.Equal(t, version, telemetry.ApplicationVersion)
	require.Len(t, telemetry.Args, 2)
	require.Len(t, telemetry.Flags, 1)
	require.Equal(t, "true", telemetry.Flags["composed"])
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

func getSummaryCmd() *cobra.Command {
	flags := DiffFlags{}

	cmd := cobra.Command{
		Use:   "summary base revision [flags]",
		Short: "Generate a diff summary",
		Long: `Display a summary of changes between base and revision specs.
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
	cmd.PersistentFlags().BoolVarP(&flags.failOnDiff, "fail-on-diff", "", false, "exit with return code 1 when any change is found")

	return &cmd
}
