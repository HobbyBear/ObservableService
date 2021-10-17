package impl

import (
	"ObservableService/services/userservice/pb"
	context "golang.org/x/net/context"
)

type Service struct {
}

func (s Service) GetUser(context context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	return &pb.GetUserResp{
		Uid:  13,
		Name: "xch",
	}, nil
}
