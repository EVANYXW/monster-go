// Package client @Author evan_yxw
// @Date 2024/8/23 16:06:00
// @Desc
package client

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"time"
)

type loginClient interface {
	GetID() uint32
	AccountResult(code int, actorID uint32, gameID uint32)
}

type Client interface {
}

type client struct {
	ID          uint32 // 客户端ID
	SID         server.ClientID
	status      ILoginStatus
	LoginResult uint32

	LoginType uint32
	LoginData []string

	acceptor *rpc.Acceptor

	closeSig chan bool

	AccountID string
	ActorID   uint32
	GameID    uint32
}

func New(id uint32, loginType uint32, loginData []string) *client {
	c := &client{
		ID: id,
	}
	server.ID2Cid(c.ID, &c.SID)

	c.acceptor = rpc.NewAcceptor(10)
	c.closeSig = make(chan bool, 1)

	c.LoginType = loginType
	c.LoginData = nil
	c.LoginData = append(c.LoginData, loginData...)

	//c.acceptor.Regist("AccountResult", c.onAccountResult)
	return c
}

func (c *client) SetStatus(statusID uint8) {
	//xsf_log.Debug("client SetStatus", xsf_log.Uint8("statusID", statusID))
	if c.status != nil {
		c.status.End(c)
		//xsf_log.Debug("client SetStatus", xsf_log.String("status end", c.status.GetName()))
	}

	c.status = manager.GetStatus(statusID)
	//xsf_log.Debug("client SetStatus", xsf_log.String("status start", c.status.GetName()))
	c.status.Start(c)
}

func (c *client) GetLoginData(id uint8) string {
	return c.LoginData[id]
}

func (c *client) GetID() uint32 {
	return c.ID
}

func (c *client) AccountResult(code int, actorID uint32, gameID uint32) {
	c.acceptor.Go("AccountResult", code, actorID, gameID)
}

func (c *client) Start() {
	async.Go(func() {
		switch c.LoginType {
		case uint32(xsf_pb.LoginType_PHXH):
			c.SetStatus(loginStatusID_PHXH)

		case uint32(xsf_pb.LoginType_TapTap):
			c.SetStatus(loginStatusID_SDKTapTap)

		case uint32(xsf_pb.LoginType_Firebase):
			c.SetStatus(loginStatusID_Firebase)

		case uint32(xsf_pb.LoginType_PHXHPhone):
			c.SetStatus(loginStatusID_PhxhPhone)
		}

	OUTLABEL:
		for {
			select {
			//case <-npCloseChan:
			//	return
			case <-c.closeSig:
				break OUTLABEL
			case callMsg, ok := <-c.acceptor.ChanCall:
				if !ok {
					break OUTLABEL
				}
				if callMsg == nil {
					continue
				}
				c.acceptor.Execute(callMsg)
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	})
}

func (c *client) GetAccountID() string {
	prefix := "phxh:"
	switch c.LoginType {
	case uint32(xsf_pb.LoginType_TapTap):
		prefix = "taptap:"
	case uint32(xsf_pb.LoginType_Firebase):
		prefix = "firebase:"
	case uint32(xsf_pb.LoginType_PHXHPhone):
		prefix = "phxhP:"
	}

	return prefix + c.AccountID
}

func (c *client) Close() {
	c.closeSig <- true
	close(c.closeSig)
}
