package export

import (
	"ObservableService/services/postservice/pb"
	"context"
	"google.golang.org/grpc"
	"log"
)

var client *grpc.ClientConn

func init() {
	var err error
	client, err = grpc.Dial("postservice:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatal("init post client fail ", err)
	}
}

func GetPost(ctx context.Context, postId string) (*pb.GetPostResp, error) {
	return pb.NewPostServiceClient(client).GetPost(ctx, &pb.GetPostReq{
		Id: postId,
	})
}
