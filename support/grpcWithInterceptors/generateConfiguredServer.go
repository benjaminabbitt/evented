package grpcWithInterceptors

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GenerateConfiguredServer(log *zap.Logger, tracer opentracing.Tracer) (grpcServer *grpc.Server) {

	grpc_zap.ReplaceGrpcLogger(log)

	opts := []grpc_zap.Option{
		grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel),
	}

	grpcServer = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			otgrpc.OpenTracingServerInterceptor(tracer),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(log, opts...),
		),
		grpc_middleware.WithStreamServerChain(
			otgrpc.OpenTracingStreamServerInterceptor(tracer),
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(log, opts...),
		))
	return grpcServer
}
