package drop

import (
	"fmt"

	"github.com/spf13/cobra"
)

var receiveCommand = &cobra.Command{
	Use:   "receive",
	Short: "Receive file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send")
	},
}
