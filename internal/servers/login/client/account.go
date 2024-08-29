// Package client @Author evan_yxw
// @Date 2024/8/26 19:26:00
// @Desc
package client

import (
	"github.com/evanyxw/monster-go/internal/db"
	"math/rand"

	//common_db "github.com/evanyxw/monster-go/internal/db"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"sync/atomic"
	"time"
)

const (
	account_data_none  = 0
	account_data_ok    = 1
	account_data_error = 2
)

var (
	rpc_account_delete       = "Delete"
	rpc_account_client_login = "ClientLogin"
	rpc_account_gm_login     = "GMLogin"
)

type account struct {
	last_active uint32 // 最后活跃时间
	acceptor    *rpc.Acceptor
	data_status atomic.Int32
	clients     []loginClient
	ai          *db.AccountInfo
	closeSig    chan bool
}

func (a *account) Init(id string) {
	a.ai = new(db.AccountInfo)
	a.ai.Id = id

	//a.acceptor.Regist(xsf_redis.RPC_NAME_RESULT, a.onRedisResult)
	//a.acceptor.Regist(xsf_mongo.RPC_NAME_RESULT, a.onMongoResult)

	a.acceptor.Regist(rpc_account_delete, a.onDelete)
	a.acceptor.Regist(rpc_account_client_login, a.onClientLogin)
	a.acceptor.Regist(rpc_account_gm_login, a.onGMLogin)

	a.data_status.Store(account_data_none)
	a.last_active = uint32(time.Now().Unix())
}

func NewAccount() *account {
	a := &account{
		closeSig: make(chan bool, 1),
		acceptor: rpc.NewAcceptor(10),
	}
	return a
}

func (a *account) Start() {
	async.Go(func() {
		defer a.Release()
		// 到redis中请求账号数据
		db.RD_GetAccountInfo(rand.Uint32(), a.ai.Id, a.acceptor)
	OUTLABEL:
		for {
			select {
			case callMsg, ok := <-a.acceptor.ChanCall: // 外部调用本地接口处理
				if !ok {
					break OUTLABEL
					//return
				}
				if callMsg == nil {
					continue
				}
				a.acceptor.Execute(callMsg)
			case <-a.closeSig:
				break OUTLABEL
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	})
}

func (a *account) Release() {
	close(a.closeSig)
	a.acceptor.Close()
}

func (a *account) End() {
	a.closeSig <- true
}

func (a *account) onClientLogin(args []interface{}) {
	client := args[0].(*client)
	a.clientLogin(client)
}

func (a *account) goClientLogin(client *client) {
	a.acceptor.Go(rpc_account_client_login, client)
}

func (a *account) clientLogin(client loginClient) {
	switch a.data_status.Load() {
	case account_data_none:
		a.clients = append(a.clients, client)
	case account_data_error:
		client.AccountResult(account_data_error, 0, 0)
	case account_data_ok:
		//client.AccountResult(account_data_ok, a.ai.ActorID, a.ai.GameID)
	}
}

func (a *account) onRedisResult(args []interface{}) {

}

func (a *account) onMongoResult(args []interface{}) {

}

func (a *account) onDelete(args []interface{}) {
	clientId := args[0].(uint32)
	a.DeleteByID(clientId)
}

func (a *account) DeleteByID(id uint32) {
	for i := 0; i < len(a.clients); i++ {
		if a.clients[i] != nil && a.clients[i].GetID() == id {
			//xsf_log.Debug("Account DeleteByID", xsf_log.Uint32("client id", id))
			a.clients[i] = nil
			break
		}
	}
}

func (a *account) onGMLogin(args []interface{}) {
	//client := args[0].(*gmClient)
	//a.clientLogin(client)
}
