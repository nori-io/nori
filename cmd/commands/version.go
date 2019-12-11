package commands

import (
	"fmt"

	"github.com/nori-io/nori/internal/version"
	"github.com/spf13/cobra"
)

var (
	// Cmd version command
	versionCmd = &cobra.Command{
		Use:           "version",
		Short:         "application version",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version.GetHumanVersion())
		},
	}
)
