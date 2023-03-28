package main

import (
	"context"
	gen "github.com/yerlanov/go-tour/grpctest/gen/go/hello"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GreeterServiceImpl struct {
	gen.UnimplementedGreeterServer
}

func (g *GreeterServiceImpl) SayHello(ctx context.Context, in *gen.HelloRequest) (*gen.HelloReply, error) {
	return &gen.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	server := grpc.NewServer()

	gen.RegisterGreeterServer(server, &GreeterServiceImpl{})

	if l, err := net.Listen("tcp", ":50551"); err != nil {
		log.Fatal("error in listening on port :50551", err)
	} else {
		// the gRPC server
		if err := server.Serve(l); err != nil {
			log.Fatal("unable to start server", err)
		}
	}
}
