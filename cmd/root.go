package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// RootCmd root命令
var RootCmd = &cobra.Command{
	Use:   "monsterGo",
	Short: "games server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}

	},
}

// Execute 执行
func Execute() {
	RootCmd.AddCommand(ServerCmd)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//defer db.CloseDB()
}
