package service

import (
	"context"

	"api_gateway/internal/client"

	pb "auth_service/proto"
	"google.golang.org/grpc/metadata"
)

type AuthService struct {
	client client.AuthClient
}

func NewAuthService(c client.AuthClient) *AuthService {
	return &AuthService{client: c}
}	

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.client.Register(ctx, req)
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.client.Login(ctx, req)
}

func (s *AuthService) SendEmail(ctx context.Context, sessionID string, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {

	md := metadata.New(map[string]string{
		"authorization": sessionID,
	})

	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.client.SendEmail(ctx, req)
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) (*pb.LogoutResponse, error) {

ctx = withAuth(ctx, sessionID)

	return s.client.Logout(ctx, &pb.LogoutRequest{})
}

func withAuth(ctx context.Context, sessionID string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": sessionID,
	})
	return metadata.NewOutgoingContext(ctx, md)
}