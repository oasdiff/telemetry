package server_test

import (
	"testing"

	"github.com/oasdiff/telemetry/server"
	"github.com/stretchr/testify/require"
)

func TestGetEventsUrl(t *testing.T) {

	require.NotEmpty(t, server.GetEventsUrl())
}
