package client

import (
	"github.com/evanyxw/game_proto/msg/messageId"
	"google.golang.org/protobuf/proto"
)

func (c *Client) InputHandlerRegister() {
	c.inputHandlers[messageId.MessageId_CSCreatePlayer.String()] = c.CreatePlayer
	c.inputHandlers[messageId.MessageId_CSLogin.String()] = c.Login
	c.inputHandlers[messageId.MessageId_CSAddFriend.String()] = c.AddFriend
	c.inputHandlers[messageId.MessageId_CSDelFriend.String()] = c.DelFriend
	c.inputHandlers[messageId.MessageId_CSSendChatMsg.String()] = c.SendChatMsg
}

func (c *Client) Transport(id messageId.MessageId, message proto.Message) {
	//bytes, err := proto.Marshal(message)
	//if err != nil {
	//	return
	//}

	//c.cli.ChMsg <- &network.Message{
	//	ID:   uint64(id),
	//	Data: bytes,
	//}
	pack, _ := c.cli.Pack(uint64(id), message)
	c.cli.SetSignal(pack) //evan
}
