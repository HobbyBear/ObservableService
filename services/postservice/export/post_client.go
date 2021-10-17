package export

import (
	"ObservableService/services/postservice/pb"
	"ObservableService/trace"
	"context"
	"google.golang.org/grpc"
	"log"
)

var client *grpc.ClientConn



func init() {
	var err error
	client ,err = grpc.Dial("127.0.0.1:8082",grpc.WithInsecure(),grpc.WithUnaryInterceptor(trace.TraceSpanClientInterceptor()))
	if err != nil{
		log.Fatal("init post client fail ",err)
	}
}

func GetPost(ctx context.Context,postId string) (*pb.GetPostResp,error) {
	return pb.NewPostServiceClient(client).GetPost(ctx,&pb.GetPostReq{
		Id:                   postId,
	})
}
