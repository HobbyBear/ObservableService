package export

import (
	"ObservableService/pkg/monitor"
	"ObservableService/services/postservice/pb"
	"context"
	"google.golang.org/grpc"
	"log"
)

var client *grpc.ClientConn

func init() {
	var err error
	client, err = grpc.Dial("postservice:8082", grpc.WithInsecure(), grpc.WithUnaryInterceptor(monitor.UnaryClientInterceptor))
	if err != nil {
		log.Fatal("init post client fail ", err)
	}
}

func GetPost(ctx context.Context, postId string) (*pb.GetPostResp, error) {
	return pb.NewPostServiceClient(client).GetPost(ctx, &pb.GetPostReq{
		Id: postId,
	})
}
