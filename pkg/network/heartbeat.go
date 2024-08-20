// Package network @Author evan_yxw
// @Date 2024/8/19 18:47:00
// @Desc
package network

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/async"
	"time"
)

func ClientHeartbeat(d time.Duration, callback func()) {
	async.Go(func() {
		//ticker := time.NewTicker(6 * time.Second)
		//defer ticker.Stop()
		//for range ticker.C {
		//	//msg := &xsf_pb.Cc_C_Heartbeat{}
		//	messageID := uint64(xsf_pb.SMSGID_Cc_C_Heartbeat)
		//	msg, _ := rpc.GetMessage(messageID)
		//	m.SendMessage(messageID, msg)
		//}

		timer := time.NewTimer(d * time.Second)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				//msg := &xsf_pb.Cc_C_Heartbeat{}
				//messageID := uint64(xsf_pb.SMSGID_Cc_C_Heartbeat)
				//msg, _ := rpc.GetMessage(messageID)
				//m.SendMessage(messageID, msg)
				callback()
			}
			_ = timer.Reset(time.Duration(d * time.Second)) //重制心跳上报时间间隔
		}
	})
}

func ServerHeartbeat(timeCount, timeOut time.Duration, lastHeart *uint64, callback func()) {
	async.Go(func() {
		timer := time.NewTimer(timeCount * time.Second)
		defer timer.Stop()

		heartbeatTimeout := uint64(time.Duration(timeOut))
	OUTLABEL:
		for {
			select {
			case <-timer.C:
				curTime := time.Now().Unix()
				if uint64(curTime) > *lastHeart+heartbeatTimeout {
					fmt.Println("ServerHeartbeat 心跳未回，被T掉")
					callback()
					break OUTLABEL
				}
			}
			_ = timer.Reset(time.Duration(timeCount * time.Second)) //重制心跳上报时间间隔
		}
	})
}
