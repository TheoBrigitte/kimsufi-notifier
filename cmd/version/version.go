package version

import (
	"fmt"

	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
)

var (
	// Version is the current version of the application
	// It is set at build time using ldflags
	Version = "n/a"
)

var (
	Cmd = &cobra.Command{
		Use:   "version",
		Short: "show version",
		RunE:  runner,
	}
)

// runner is the main function for the version command
func runner(cmd *cobra.Command, args []string) error {
	fmt.Println(version.Print("kimsufi-notifier"))

	return nil
}
