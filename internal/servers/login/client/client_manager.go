// Package client @Author evan_yxw
// @Date 2024/8/23 19:18:00
// @Desc
package client

import (
	"github.com/evanyxw/monster-go/pkg/rpc"
	"time"
)

var manager *Manager

type Manager struct {
	status         []ILoginStatus
	rpcAcceptor    *rpc.Acceptor
	clients        map[uint32]*client
	accounts       map[string]*account
	login_accounts map[uint32]*account

	accountCreate    uint64
	clientCreate     uint64
	clientDelete     uint64
	accountCheckTime uint32
}

func NewClientManager() *Manager {
	r := &Manager{
		accounts:       make(map[string]*account),
		login_accounts: make(map[uint32]*account),
		clients:        make(map[uint32]*client),
	}
	r.CreateStatus()
	manager = r
	return r
}

func (m *Manager) Init(rpcAcceptor *rpc.Acceptor) {
	m.rpcAcceptor = rpcAcceptor
	m.RegisterRPC()
}

func (m *Manager) RegisterRPC() {
	m.rpcAcceptor.Regist("NewAccount", m.onNewAccount)
}

func (m *Manager) NewClient(id uint32, loginType uint32, loginData []string) *client {
	_, ok := m.clients[id]
	if ok {
		//xsf_log.Warn("loginManager GetClient client exist", xsf_log.Uint32("id", id))
		return nil
	} else {
		c := New(id, loginType, loginData)
		m.clientCreate++
		m.clients[id] = c
		return c
	}
}

func (m *Manager) GetAccount(id string) *account {
	a, ok := m.accounts[id]
	if ok {
		//xsf_log.Debug("find exist account", xsf_log.String("id", id))
	} else {
		a = NewAccount()
		a.Init(id)
		m.accountCreate++
		m.accounts[id] = a
		a.Start()
	}

	return a
}

func (m *Manager) onNewAccount(args []interface{}) {
	client := args[0].(*client)
	id := args[1].(string)
	account := m.GetAccount(id)
	account.last_active = uint32(time.Now().Unix())
	account.goClientLogin(client)

	m.login_accounts[client.ID] = account
}

func (m *Manager) NewAccount(c *client, id string) {
	m.rpcAcceptor.Go("NewAccount", c, id)
}

func (m *Manager) CreateStatus() {
	m.status = make([]ILoginStatus, loginStatusID_Max)

	m.status[loginStatusID_PHXH] = new(loginStatus_PHXH)
	//m.status[loginStatusID_SDKTapTap] = new(loginStatus_TapTap)
	m.status[loginStatusID_Account] = new(loginStatus_Account)
	//m.status[loginStatusID_Login2Game] = new(loginStatus_Login2Game)
	//m.status[loginStatusID_Error] = new(loginStatus_Error)
	//m.status[loginStatusID_Firebase] = new(loginStatus_Firebase)
	//m.status[loginStatusID_PhxhPhone] = new(loginStatus_PhxhPhone)
}

func (m *Manager) GetStatus(status uint8) ILoginStatus {
	return m.status[status]
}

func (m *Manager) CloseClient(clientID uint32) {
	//clientID := args[0].(uint32)
	//client, ok := lm.clients[clientID]
	//account, ok2 := lm.login_accounts[clientID]
	//
	//if ok2 {
	//	account.GoDeleteClient(clientID)
	//	delete(lm.login_accounts, clientID)
	//}
	//
	//if ok {
	//	delete(lm.clients, clientID)
	//	client.End()
	//	lm.clientDelete++
	//}

	client, ok := m.clients[clientID]
	if ok {
		delete(m.clients, clientID)
		client.Close()
		m.clientDelete++
	}

	account, ok := m.login_accounts[clientID]
	if ok {
		account.GoDeleteClient(clientID)
		delete(m.login_accounts, clientID)
	}
}

func (m *Manager) OnUpdate() {
	current := uint32(time.Now().Unix())
	if current > m.accountCheckTime+60 {
		m.accountCheckTime = current

		lifeMap := make(map[string]*account)

		for id, a := range m.accounts {
			if a.data_status.Load() == account_data_error {
				a.End()
			} else {
				if current > a.last_active+120 {
					a.End()
				} else {
					lifeMap[id] = a
				}
			}
		}

		m.accounts = lifeMap
	}
}
