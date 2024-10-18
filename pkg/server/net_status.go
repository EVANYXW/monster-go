// Package server @Author evan_yxw
// @Date 2024/10/16 19:14:00
// @Desc
package server

// 为服务网络节点的状态，为了表示已经准备好可以其他服务器来链接
const (
	Net_RunStep_Start = iota + 1
	Net_RunStep_Done
)

func NetStatusStart(status *int) {
	*status = Net_RunStep_Start
}

func NetStatusDone(status *int) {
	*status = Net_RunStep_Done
}

func NetStatusIsDone(status int) bool {
	return status == Net_RunStep_Done
}
