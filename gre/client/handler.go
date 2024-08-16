package client

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/network"
)

type MessageHandler func(packet *network.Packet)

type InputHandler func(param *InputParam)

func (c *Client) Test(param *InputParam) {
	//id := c.cli.TcpConn.GetMessageIdByCmd(param.Command)
	if len(param.Param) != 1 {
		return
	}

	pb := &xsf_pb.Clt_L_Login{}
	c.cli.SendMessage(pb)
	//msg := &xsf_pb.Cc_C_Handshake{}
	//c.Transport(id, msg)

}
