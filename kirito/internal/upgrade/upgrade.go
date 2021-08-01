package upgrade

import (
	"fmt"

	"github.com/go-kirito/cmd/kirito/internal/base"

	"github.com/spf13/cobra"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the kirito tools",
	Long:  "Upgrade the kirito tools. Example: kirito upgrade",
	Run:   Run,
}

// Run upgrade the kratos tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoGet(
		"github.com/go-kirito/cmd/kirito",
		"github.com/go-kirito/cmd/protoc-gen-go-kirito",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"github.com/envoyproxy/protoc-gen-validate",
	)
	if err != nil {
		fmt.Println(err)
	}
}
