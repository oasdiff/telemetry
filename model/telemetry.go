package model

import (
	"os"
	"runtime"
	"time"

	machine "github.com/denisbrodbeck/machineid"
	"golang.org/x/exp/slog"
)

const (
	Application    = "oasdiff"
	EnvNoTelemetry = "OASDIFF_NO_TELEMETRY"
	DefaultTimeout = time.Millisecond * 700
)

type Telemetry struct {
	Id                 string
	Time               int64
	MachineId          string
	Runtime            string // darwin/windows
	Platform           string // docker/github-action
	Command            string
	Args               []string
	Flags              map[string]string
	Application        string
	ApplicationVersion string
	// Duration           int64
}

func NewTelemetry(app string, appVersion string, cmd string, args []string, flags map[string]string) *Telemetry {

	machineId, err := machine.ProtectedID(app)
	if err != nil {
		slog.Debug("failed to get machine ID", "error", err)
		machineId = "na"
	}

	return &Telemetry{
		Time:               time.Now().UnixMilli(),
		Application:        app,
		ApplicationVersion: appVersion,
		MachineId:          machineId,
		Runtime:            runtime.GOOS,
		Platform:           getPlatform(),
		Command:            cmd,
		Args:               args,
		Flags:              flags,
	}
}

func getPlatform() string {

	if res := os.Getenv("PLATFORM"); res != "" {
		return res
	}

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "dockerenv"
	}

	return "na"
}
