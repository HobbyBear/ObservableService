package impl

import (
	"ObservableService/services/postservice/pb"
	"ObservableService/services/userservice/export"
	context "golang.org/x/net/context"
	"log"
)

type Service struct {

}

func (s Service) GetPost(context context.Context, req *pb.GetPostReq) (*pb.GetPostResp, error) {
	userResp,err := export.GetUser(context,13)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	return &pb.GetPostResp{
		Uid:                  userResp.Uid,
		Text:                 "get post success",
	},nil
}

