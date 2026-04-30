package grpc

import (
	"auth_service/pkg/logger"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		logger.Infof("Incoming RPC: %s", info.FullMethod)

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			logger.Infof("Metadata received: %v", md)

			values := md.Get("authorization")
			if len(values) > 0 {
				sessionID := values[0]
				logger.Infof("Extracted session: %s", sessionID)

				ctx = context.WithValue(ctx, SessionKey, sessionID)
			} else {
				logger.Warn("No authorization header found")
			}
		} else {
			logger.Warn("No metadata found in context")
		}

		return handler(ctx, req)
	}
}