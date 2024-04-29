package drop

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sendCommand = &cobra.Command{
	Use:   "send",
	Short: "Send file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send")
	},
}
