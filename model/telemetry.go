package model

import (
	"os"
	"runtime"
	"time"

	machine "github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slog"
)

const (
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
	ApplicationVersion string
	// Duration           int64
}

func FromCommand(cmd *cobra.Command) *Telemetry {

	subCommandName := ""
	args := []string{}
	flagNameToValue := make(map[string]string)

	for _, currSubCommand := range cmd.Commands() {
		subCommandName = currSubCommand.CalledAs()
		if subCommandName != "" {
			currSubCommand.Flags().Visit(func(flag *pflag.Flag) {
				flagNameToValue[flag.Name] = flag.Value.String()
			})
			for _, currArg := range currSubCommand.Flags().Args() {
				args = append(args, currArg)
			}
			break
		}
	}

	return newTelemetry(cmd.Version, subCommandName, args, flagNameToValue)
}

func newTelemetry(appVersion string, cmd string, args []string, flags map[string]string) *Telemetry {

	machineId, err := machine.ID()
	if err != nil {
		slog.Debug("failed to get machine ID", "error", err)
		machineId = "na"
	}

	return &Telemetry{
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

	if res := os.Getenv("PLATFORM"); res != "" {
		return res
	}

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "dockerenv"
	}

	return "na"
}
