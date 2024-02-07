package cmd

import (
	"bilibili/monster-go/configs"
	"bilibili/monster-go/internal/server"
	"bilibili/monster-go/internal/server/factory"
	"bilibili/monster-go/internal/server/world"
	"bilibili/monster-go/pkg/env"
	"github.com/spf13/cobra"
	"log"
	_ "net/http/pprof" // for side effects only
)

var (
	serverName string
	envStr     string
)

func init() {
	// 游戏服务注册
	factory.Register(server.World, world.New)

	// 启动服务参数
	ServerCmd.Flags().StringVar(&envStr, "env", "", "env")
	ServerCmd.Flags().StringVar(&serverName, "server_name", "", "server_name")

}

// ServerCmd server 服务的cmd方法、
var ServerCmd = &cobra.Command{
	Use:   "run",
	Short: "Run game server",
	Run: func(cmd *cobra.Command, args []string) {
		env.Init(envStr)
		configs.Init()

		if serverName == "" {
			log.Fatal("Please specify a server name")
		}

		Run(serverName)
	},
}
