package drop

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop is a simple CLI tool for file sharing.",
	Long: `Drop is a simple CLI tool for file sharing.
It allows you to send files to other devices on the same network.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

// Execute runs the root command.
func Execute() {
	rootCmd.AddCommand(sendCommand)
	rootCmd.AddCommand(receiveCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
