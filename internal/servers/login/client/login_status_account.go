// Package client @Author evan_yxw
// @Date 2024/8/23 18:57:00
// @Desc
package client

type loginStatus_Account struct {
}

func (ls *loginStatus_Account) GetID() uint8 {
	return loginStatusID_Account
}

func (ls *loginStatus_Account) GetName() string {
	return "loginStatus_Account"
}

func (ls *loginStatus_Account) Start(client *client) {
	manager.NewAccount(client, client.GetAccountID())
}

func (ls *loginStatus_Account) End(client *client) {

}

func (ls *loginStatus_Account) OnAccountResult(client *client, code int) {
}
