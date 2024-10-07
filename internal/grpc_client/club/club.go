// Package login @Author evan_yxw
// @Date 2024/10/5 21:40:00
// @Desc
package club

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/pkg/grpcpool"
	"github.com/evanyxw/monster-go/proto/pb"
)

func NewClubServiceClient() (pb.ClubSerClient, error) {
	conn := grpcpool.Instance().Dial(servers.LoginGrpc)
	return pb.NewClubSerClient(conn), nil
}
