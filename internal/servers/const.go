package servers

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

var (
	NodeManager     module.NodeManager
	NetPointManager network.INPManager
	ConnectorKernel *module.ConnectorKernel
)

const (
	Rpc = "_rpc"
)

const (
	Logic  = "logic"
	World  = "world"
	Center = "center"
	Gate   = "gate"
)
