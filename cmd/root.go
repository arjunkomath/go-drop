package drop

import (
	"drop/pkg/receive"
	"drop/pkg/send"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1-beta"

var rootCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop is a simple CLI tool for file sharing.",
	Long: `Drop is a simple CLI tool for file sharing.
It allows you to send files to other devices on the same network.`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute runs the root command.
func Execute() {
	rootCmd.AddCommand(send.Command)
	rootCmd.AddCommand(receive.Command)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
