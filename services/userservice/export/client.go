package export

import (
	"ObservableService/pkg/monitor"
	"ObservableService/services/userservice/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var client *grpc.ClientConn

func init() {
	var err error
	client, err = grpc.Dial("userservice:8083", grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(monitor.TraceUnaryClientInterceptor(), monitor.UnaryMetricClientInterceptor))
	if err != nil {
		log.Fatal("init post client fail ", err)
	}
}

func GetUser(ctx context.Context, uid int) (*pb.GetUserResp, error) {
	return pb.NewUserServiceClient(client).GetUser(ctx, &pb.GetUserReq{
		Id: int64(uid),
	})
}
