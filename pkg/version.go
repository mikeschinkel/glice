package glice

import (
	"fmt"
	"github.com/spf13/cobra"
)

func RunVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("%s v%s\n", AppName, AppVersion)
}
