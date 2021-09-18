package main

import (
	"log"

	"github.com/go-kirito/cmd/kirito/internal/change"
	"github.com/go-kirito/cmd/kirito/internal/domain"
	"github.com/go-kirito/cmd/kirito/internal/model"
	"github.com/go-kirito/cmd/kirito/internal/project"
	"github.com/go-kirito/cmd/kirito/internal/proto"
	"github.com/go-kirito/cmd/kirito/internal/run"
	"github.com/go-kirito/cmd/kirito/internal/upgrade"
	"github.com/go-kirito/cmd/kirito/internal/wire"

	"github.com/spf13/cobra"
)

var (
	version string = "v1.0.0"

	rootCmd = &cobra.Command{
		Use:     "kirito",
		Short:   "kirito: An elegant toolkit for Go microservices.",
		Long:    `kirito: An elegant toolkit for Go microservices.`,
		Version: version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(change.CmdChange)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(wire.CmdWire)
	rootCmd.AddCommand(domain.CmdDomain)
	rootCmd.AddCommand(model.CmdModel)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
