package cmd

import (
	"bilibili/monster-go/configs"
	"bilibili/monster-go/internal/server/factory"
	"bilibili/monster-go/internal/server/world"
	"bilibili/monster-go/pkg/env"
	"fmt"
	"github.com/spf13/cobra"
	_ "net/http/pprof" // for side effects only
)

var (
	serverName string
	envStr     string
)

func init() {
	// 游戏服务注册
	factory.Register("world", world.New)

	// 启动服务参数
	ServerCmd.Flags().StringVar(&envStr, "env", "", "env")
	ServerCmd.Flags().StringVar(&serverName, "server_name", "", "server_name")

}

// ServerCmd server 服务的cmd方法、
var ServerCmd = &cobra.Command{
	Use:   "run",
	Short: "games world server",
	Run: func(cmd *cobra.Command, args []string) {
		env.Init(envStr)
		configs.Init()

		if serverName == "" {
			fmt.Println("Please specify a server name")
			return
		}
		Run(serverName)
	},
}
