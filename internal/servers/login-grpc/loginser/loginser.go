// Package loginser @Author evan_yxw
// @Date 2024/10/4 11:32:00
// @Desc
package loginser

import (
	"context"
	"github.com/evanyxw/monster-go/proto/pb"
	"google.golang.org/grpc"
)

type Login struct {
	pb.UnimplementedClubSerServer
}

func NewLogin() *Login {
	return &Login{}
}

func (l *Login) ClubCreate(ctx context.Context, req *pb.ClubCreate_Req) (*pb.ClubCreate_Rsp, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Login) ClubRead(ctx context.Context, req *pb.ClubRead_Req) (*pb.ClubRead_Rsp, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Login) ClubDelete(ctx context.Context, req *pb.ClubDelete_Req) (*pb.ClubDelete_Rsp, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Login) ClubRoleInfo(ctx context.Context, req *pb.ClubRoleInfo_Req) (*pb.ClubRoleInfo_Rsp, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Login) ClubJoin(ctx context.Context, req *pb.ClubJoin_Req) (*pb.ClubJoin_Rsp, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Login) ClubExit(ctx context.Context, req *pb.ClubExit_Req) (*pb.ClubExit_Rsp, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Login) mustEmbedUnimplementedClubSerServer() {
	//TODO implement me
	panic("implement me")
}

func (l *Login) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterClubSerServer(t, l)
		return nil
	}
}
