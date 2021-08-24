package main

import (
	"github.com/go-kirito/cmd/kirito/internal/ioc"
	"log"

	"github.com/go-kirito/cmd/kirito/internal/change"
	"github.com/go-kirito/cmd/kirito/internal/project"
	"github.com/go-kirito/cmd/kirito/internal/proto"
	"github.com/go-kirito/cmd/kirito/internal/run"
	"github.com/go-kirito/cmd/kirito/internal/upgrade"

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
	rootCmd.AddCommand(ioc.CmdIoc)
	rootCmd.AddCommand(run.CmdRun)
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
