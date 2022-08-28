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
	err := base.GoInstall(
		"github.com/go-kirito/cmd/kirito@latest",
		"github.com/go-kirito/cmd/protoc-gen-go-kirito@latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"github.com/envoyproxy/protoc-gen-validate@latest",
		"github.com/google/wire/cmd/wire",
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
