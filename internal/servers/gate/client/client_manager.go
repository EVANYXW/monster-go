// Package client @Author evan_yxw
// @Date 2024/8/7 22:42:00
// @Desc
package client

import (
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
)

type clientManager struct {
	clients   []*Client
	GlobalKey uint8
}

func NewClientManager() *clientManager {
	c := &clientManager{}
	c.Init()
	return c
}

func (c *clientManager) Init() {
	c.clients = make([]*Client, 5000)

}

func (c *clientManager) NewClient(np *network.NetPoint) (module_def.Client, bool) {
	for i := 0; i < len(c.clients); i++ {
		if c.clients[i] == nil {
			client := NewClient(np)
			client.Init()
			c.clients[i] = client
			client.CID.Gate = server.SID.Index
			client.CID.ID = uint16(i)
			client.CID.Key = c.GlobalKey
			//xsf_log.Debug("NewClient new", xsf_log.Uint8("gate", client.CID.Gate), xsf_log.Uint16("id", client.CID.ID), xsf_log.Uint8("key", client.CID.Key))
			client.ID.Store(server.Cid2ID(&client.CID))
			c.GlobalKey++
			return client, true
		} else if c.clients[i].ID.Load() == 0 {
			client := c.clients[i]
			client.netPoint = np
			client.Init()
			client.CID.Gate = server.SID.Index
			client.CID.ID = uint16(i)
			client.CID.Key = c.GlobalKey
			c.GlobalKey++
			client.ID.Store(server.Cid2ID(&client.CID))
			//client := ck.clients[i]
			//client.CID.Gate = xsf_server.SID.Index
			//client.CID.ID = uint16(i)
			//client.CID.Key = ck.GlobalKey
			//ck.GlobalKey++
			//client.ID.Set(xsf_server.Cid2ID(&client.CID))
			//
			//client.reInit(conn)
			//
			//xsf_log.Debug("NewClient exist", xsf_log.Uint8("gate", client.CID.Gate), xsf_log.Uint16("id", client.CID.ID), xsf_log.Uint8("key", client.CID.Key))
			return client, false
		}
	}
	return nil, false
}

func (c *clientManager) GetClient(id uint32) module_def.Client {
	var CID server.ClientID
	server.ID2Cid(id, &CID)

	if c.clients[CID.ID] != nil && c.clients[CID.ID].ID.Load() == id {
		if c.clients[CID.ID].isHandshake.Load() {
			return c.clients[CID.ID]
		} else {
			return nil
		}
	} else {
		return nil
	}
}
