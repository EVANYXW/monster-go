package cmd

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/center"
	"github.com/evanyxw/monster-go/internal/servers/gate"
	gate_grpc "github.com/evanyxw/monster-go/internal/servers/gate-grpc"
	"github.com/evanyxw/monster-go/internal/servers/login"
	cmdapk "github.com/evanyxw/monster-go/pkg/cmd"
	"github.com/evanyxw/monster-go/pkg/server/engine"
	"log"
	_ "net/http/pprof" // for side effects only

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
	engine.Register(servers.GateGrpc, gate_grpc.New)

	// 启动服务参数
	ServerCmd.Flags().StringVar(&envStr, "env", "", "env")
	ServerCmd.Flags().StringVar(&serverName, "server_name", "", "server_name")

}

// ServerCmd server 服务的cmd方法、
var ServerCmd = &cobra.Command{
	Use:   "run",
	Short: "run game server",
	Run: func(cmd *cobra.Command, args []string) {
		if serverName == "" {
			log.Fatal("Please specify a server name")
		}

		env.SetActive(envStr)

		//_, _ = logger.NewJSONLogger(
		//	logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		//	logger.WithTimeLayout(timeutil.CSTLayout),
		//	logger.WithFileP(configs.LogFile, serverName),
		//)
		//zap_log.NewLogger()

		cmdapk.Run(serverName)
	},
}
