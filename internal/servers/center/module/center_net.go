package module

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/servers/center/handler"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
)

type CenterNet struct {
	kernel       module.IModuleKernel
	status       int
	startIndex   int
	curStartNode *configs.ServerNode
}

func NewCenterNet(id int32, maxConnNum uint32, info server.Info) *CenterNet {
	centerCnf := configs.Get().Center
	info.Ip = centerCnf.Ip
	info.Port = centerCnf.Port

	centerNet := &CenterNet{
		kernel: module.NewNetKernel(
			maxConnNum,
			info,
			handler.NewCenterNetMsg(),
			new(network.DefaultPackerFactory),
			module.WithNoWaitStart(true)),
	}

	module.NewBaseModule(id, centerNet)
	network.NetPointManager = centerNet.kernel.GetNPManager()
	return centerNet
}

func (c *CenterNet) Init(baseModule *module.BaseModule) bool {
	c.kernel.Init(baseModule)
	return true
}

func (c *CenterNet) DoRegister() {
	c.kernel.DoRegister()
}

func (c *CenterNet) DoRun() {
	c.kernel.DoRun()
	c.status = server.CN_RunStep_StartServer
	c.startIndex = 0
}

func (c *CenterNet) DoWaitStart() {

}

func (c *CenterNet) DoRelease() {
	c.kernel.DoRelease()
}

func (c *CenterNet) OnOk() {

}

func (c *CenterNet) OnStartCheck() int {
	serverCnf := configs.Get()
	if !serverCnf.AutoStart {
		return module.ModuleRunCode_Ok
	}

	serverList := configs.Get().ServerList
	switch c.status {
	case server.CN_RunStep_StartServer:
		c.curStartNode = &(serverList[c.startIndex])
		dir, _ := os.Getwd()

		// 兼容开发时的直接运行
		binDir := dir + "/bin"
		cmdStr := "./bin/nld_server run --server_name " + c.curStartNode.EPName

		_, err := os.Stat(binDir)
		if os.IsNotExist(err) {
			logger.Info("找不到bin文件夹，执行当前目录sh文件")
			cmdStr = "./single_start.sh " + c.curStartNode.EPName
		}

		cmdFields := strings.Fields(cmdStr)
		cmd := exec.Command(cmdFields[0], cmdFields[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		if err != nil {
			logger.Error("Error running command:", zap.Error(err))
		}

		// fixMe 恢复
		c.status = server.CN_RunStep_WaitHandshake
		//c.status = server.CN_RunStep_HandshakeDone
	case server.CN_RunStep_HandshakeDone:
		c.startIndex++
		if c.startIndex >= len(serverList) {
			return module.ModuleRunCode_Ok
		} else {
			c.status = server.CN_RunStep_StartServer
		}
	}

	return module.ModuleRunCode_Wait
}

func (c *CenterNet) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *CenterNet) Update() {

}

func (c *CenterNet) GetKernel() module.IModuleKernel {
	return c.kernel
}
