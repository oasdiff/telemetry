package model

import (
	"os"
	"runtime"
	"time"

	machine "github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
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
	ApplicationVersion string
	// Duration           int64
}

func NewTelemetry(appVersion string, cmd string, args []string, flags map[string]string) Telemetry {

	machineId, err := machine.ID()
	if err != nil {
		slog.Debug("failed to get machine ID", "error", err)
		machineId = "na"
	}

	return Telemetry{
		Id:                 uuid.NewString(),
		Time:               time.Now().UnixMilli(),
		MachineId:          machineId,
		ApplicationVersion: appVersion,
		Runtime:            runtime.GOOS,
		Platform:           getPlatform(),
		Command:            cmd,
		Args:               args,
		Flags:              flags,
	}
}

func getPlatform() string {

	if res := os.Getenv("Platform"); res != "" {
		return res
	}

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "dockerenv"
	}

	return "na"
}
