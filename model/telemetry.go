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
	KeyEvents      = "events"
	EnvNoTelemetry = "OASDIFF_NO_TELEMETRY"
	DefaultTimeout = time.Millisecond * 700
)

type Telemetry struct {
	Id                 string            `json:"id"`
	Time               int64             `json:"time"`
	MachineId          string            `json:"machine_id"`
	Runtime            string            `json:"runtime"`  // darwin/windows
	Platform           string            `json:"platform"` // docker/github-action
	Command            string            `json:"command"`
	Args               []string          `json:"args"`
	Flags              map[string]string `json:"flags"`
	Application        string            `json:"application"`
	ApplicationVersion string            `json:"application_version"`
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
