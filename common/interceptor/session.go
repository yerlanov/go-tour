package interceptor

import (
	"context"
	"encoding/json"
	"github.com/yerlanov/go-tour/common/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func SessionInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if sessionJson, ok := md["session"]; ok {
		var sess session.Content
		err := json.Unmarshal([]byte(sessionJson[0]), &sess)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid session")
		}

		ctx = context.WithValue(ctx, "session", sess)
	}

	return handler(ctx, req)
}
