package export

import (
	"ObservableService/services/userservice/pb"
	"ObservableService/trace"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var client *grpc.ClientConn

func init() {
	var err error
	client, err = grpc.Dial("127.0.0.1:8083",grpc.WithInsecure(),grpc.WithUnaryInterceptor(trace.TraceSpanClientInterceptor()))
	if err != nil {
		log.Fatal("init post client fail ", err)
	}
}

func GetUser(ctx context.Context, uid int) (*pb.GetUserResp, error) {
	return pb.NewUserServiceClient(client).GetUser(ctx, &pb.GetUserReq{
		Id: int64(uid),
	})
}
