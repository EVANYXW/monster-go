package center

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/kernel"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery/center/handler"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
)

type CenterNet struct {
	kernel       module_def.IKernel
	status       int
	startIndex   int
	curStartNode *configs.ServerNode
	id           int32
}

func NewCenterNet(id int32, maxConnNum uint32) *CenterNet {
	centerCnf := configs.All().Center
	server.SetInfoIP(centerCnf.Ip)
	server.SetInfoPort(centerCnf.Port)

	centerNet := &CenterNet{
		id: id,
		kernel: kernel.NewNetKernel(
			maxConnNum,
			handler.NewCenterNetMsg(),
			new(network.DefaultPackerFactory),
			kernel.WithNoWaitStart(true)),
	}

	//module.NewBaseModule(id, centerNet) 我在试一试
	network.NetPointManager = centerNet.kernel.GetNPManager()

	return centerNet
}

func (c *CenterNet) Init(baseModule module_def.IBaseModule) bool {
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
	serverCnf := configs.All()
	if !serverCnf.AutoStart {
		return module_def.ModuleOk()
	}

	serverList := configs.All().ServerList
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
			return module_def.ModuleOk()
		} else {
			c.status = server.CN_RunStep_StartServer
		}
	}

	return module_def.ModuleWait()
}

func (c *CenterNet) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *CenterNet) Update() {

}

func (c *CenterNet) GetID() int32 {
	return c.id
}

func (c *CenterNet) GetKernel() module_def.IKernel {
	return c.kernel
}
