package world

import (
	"fmt"
	"github.com/evanyxw/game_proto/msg/player"
	"github.com/evanyxw/monster-go/pkg/network"
	"google.golang.org/protobuf/proto"
)

func (w *World) CreatePlayer(message *network.Packet) {

	msg := &player.CSLogin{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		fmt.Println(err)
	}

	//res := &player.SCCreateUser{}
	fmt.Println("准备获取id了:SCCreatePlayer")
	id := message.Conn.GetMessageIdByCmd("SCCreatePlayer")
	fmt.Println("message_id:", id)
	pack, _ := message.Conn.Pack(uint64(id), &player.SCCreateUser{})

	err = message.Conn.PackWrite(pack)
}

func (w *World) UserLogin(message *network.Packet) {
	fmt.Println("登录成功")
}
