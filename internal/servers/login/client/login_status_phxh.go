// Package client @Author evan_yxw
// @Date 2024/8/23 18:57:00
// @Desc
package client

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"unicode"
)

type loginStatus_PHXH struct {
}

func (ls *loginStatus_PHXH) GetID() uint8 {
	return loginStatusID_PHXH
}

func (ls *loginStatus_PHXH) GetName() string {
	return "loginStatus_PHXH"
}

func (ls *loginStatus_PHXH) Start(client *client) {
	account := client.GetLoginData(uint8(xsf_pb.PHXHLoginData_Account))
	if len(account) <= 0 {
		client.LoginResult = uint32(xsf_pb.LoginResult_AccountInvalid)
		client.SetStatus(loginStatusID_Error)
		return
	}

	for _, char := range account {
		if char > 127 {
			client.LoginResult = uint32(xsf_pb.LoginResult_AccountInvalid)
			client.SetStatus(loginStatusID_Error)
			return
		}

		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			client.LoginResult = uint32(xsf_pb.LoginResult_AccountInvalid)
			client.SetStatus(loginStatusID_Error)
			return
		}
	}

	//if client.GetLoginData(uint8(xsf_pb.PHXHLoginData_Password)) != LoginPass {
	//	client.LoginResult = uint32(xsf_pb.LoginResult_LoginAuthError)
	//	client.SetStatus(loginStatusID_Error)
	//	return
	//}

	client.AccountID = client.GetLoginData(uint8(xsf_pb.PHXHLoginData_Account))
	client.SetStatus(loginStatusID_Account)
}

func (ls *loginStatus_PHXH) End(client *client) {

}

func (ls *loginStatus_PHXH) OnAccountResult(client *client, code int) {
	//xsf_log.Error("loginStatus_PHXH OnAccountResult error call")
}
