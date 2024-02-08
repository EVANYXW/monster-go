package cmd

import (
	"log"
	_ "net/http/pprof" // for side effects only

	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/server"
	"github.com/evanyxw/monster-go/internal/server/factory"
	"github.com/evanyxw/monster-go/internal/server/world"
	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/spf13/cobra"
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
