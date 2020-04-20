package grpcZap

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

func GenerateConfiguredConn(target string, log *zap.SugaredLogger) *grpc.ClientConn {

	var customFunc grpc_zap.CodeToLevel
	zapLogOpts := []grpc_zap.Option{
		grpc_zap.WithLevels(customFunc),
	}
	grpc_zap.ReplaceGrpcLogger(log.Desugar())

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted, codes.Unavailable, codes.Unimplemented, codes.Unknown),
	}
	conn, err := grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(grpc_zap.UnaryClientInterceptor(log.Desugar(), zapLogOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithStreamInterceptor(grpc_zap.StreamClientInterceptor(log.Desugar(), zapLogOpts...)),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
	)
	if err != nil {
		log.Error(err)
	}
	return conn
}
