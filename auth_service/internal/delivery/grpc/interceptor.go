package grpc

import (
	"auth_service/pkg/jwt"
	"auth_service/pkg/logger"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
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

			authHeaders := md.Get("authorization")
			if len(authHeaders) > 0 && authHeaders[0] != "" {
				token := authHeaders[0]
				if strings.HasPrefix(token, "Bearer ") {
					token = strings.TrimPrefix(token, "Bearer ")
				}
				userID, err := jwt.ParseToken(token)
				if err != nil {
					return nil, status.Error(codes.Unauthenticated, "invalid token")
				}

				ctx = context.WithValue(ctx, SessionKey, userID)
				ctx = context.WithValue(ctx, TokenKey, token)
			} else {
				logger.Warn("No authorization header found")
			}
		} else {
			logger.Warn("No metadata found in context")
		}

		return handler(ctx, req)
	}
}