// Package client @Author evan_yxw
// @Date 2024/8/23 18:17:00
// @Desc
package client

const (
	loginStatusID_PHXH       = 0 // 自定义登录
	loginStatusID_SDKTapTap  = 1 // taptap登录
	loginStatusID_Account    = 2 // 获取账号数据
	loginStatusID_Login2Game = 3 // 执行game登录
	loginStatusID_Error      = 4
	loginStatusID_Firebase   = 5 // firebase登录
	loginStatusID_PhxhPhone  = 6 // phxh手机验证码登录

	loginStatusID_Max = 7
)

type ILoginStatus interface {
	GetID() uint8
	GetName() string
	Start(client *client)
	End(client *client)
	OnAccountResult(client *client, code int)
}
