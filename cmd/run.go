package cmd

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/center"
	"github.com/evanyxw/monster-go/internal/servers/gate"
	"github.com/evanyxw/monster-go/internal/servers/login"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/server/engine"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	"log"
	_ "net/http/pprof" // for side effects only

	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/spf13/cobra"
)

var (
	serverName string
	envStr     string
)

func init() {
	// 游戏服务注册
	engine.Register(servers.Gate, gate.New)
	engine.Register(servers.Center, center.New)
	engine.Register(servers.Login, login.New)

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

		logger.NewLogger(
			logger.WithDisableConsole(),
			logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
			logger.WithTimeLayout(timeutil.CSTLayout),
			logger.WithFileP(configs.LogFile, serverName),
		)

		Run(serverName)
	},
}
