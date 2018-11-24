package command

import (
	"github.com/PacketFire/immigrant/command/version"
	"github.com/Packetfire/immigrant/command/converge"

	"github.com/mitchellh/cli"
)

const (
	cliVersion string = "0.0.1"
	coutput    string = "converge command"
)

func init() {
	Register("version", func(ui cli.Ui) (cli.Command, error) { return version.New(ui, cliVersion), nil })
	Register("converge", func(ui cli.Ui) (cli.Command, error) { return converge.New(ui, coutput), nil })
}
