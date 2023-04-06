package main

import (
	"context"
	"fmt"
	"github.com/yerlanov/go-tour/common/interceptor"
	"github.com/yerlanov/go-tour/common/session"
	gen "github.com/yerlanov/go-tour/grpctest2/gen/go/calculator"
	"google.golang.org/grpc"
	"log"
	"net"
)

type CalculatorServiceImpl struct {
	gen.UnimplementedCalculatorServer
}

func (c *CalculatorServiceImpl) Multiply(ctx context.Context, in *gen.HelloRequest) (*gen.HelloReply, error) {
	sess, err := session.GetSessionFromContext(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(sess)

	result := fmt.Sprintf("Result: %d", in.Nu*in.Nu2)
	return &gen.HelloReply{Message: result}, nil
}

func main() {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.SessionInterceptor),
	)

	gen.RegisterCalculatorServer(server, &CalculatorServiceImpl{})

	if l, err := net.Listen("tcp", ":50552"); err != nil {
		log.Fatal("error in listening on port :50552", err)
	} else {
		// the gRPC server
		if err := server.Serve(l); err != nil {
			log.Fatal("unable to start server", err)
		}
	}
}
