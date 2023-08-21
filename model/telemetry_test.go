package model_test

import (
	"testing"

	"github.com/oasdiff/telemetry/model"
	"github.com/stretchr/testify/require"
)

func TestNewTelemetry(t *testing.T) {

	const cmd, version = "breaking", "v1.2.3"

	telemetry := model.NewTelemetry(version, cmd,
		[]string{"v1.yaml", "v2.yaml"},
		map[string]string{"composed": "true", "format": "text"})

	require.NotEmpty(t, telemetry.Id)
	require.True(t, telemetry.Time > 0)
	require.NotEmpty(t, telemetry.MachineId)
	require.NotEmpty(t, telemetry.Runtime)
	require.NotEmpty(t, telemetry.Platform)
	require.Equal(t, cmd, telemetry.Command)
	require.Equal(t, version, telemetry.ApplicationVersion)
	require.Len(t, telemetry.Args, 2)
	require.Len(t, telemetry.Flags, 2)
}
