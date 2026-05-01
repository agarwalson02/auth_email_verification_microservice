package grpc

import (
	"context"
	"errors"

	"auth_service/internal/models"
	"auth_service/internal/services"
	pb "auth_service/proto"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	service *services.AuthService
}

func NewHandler(service *services.AuthService) *Handler {
	return &Handler{service: service}
}

func extractSessionFromContext(ctx context.Context) string {
	val := ctx.Value(SessionKey)
	if val == nil {
		return ""
	}

	return val.(string)
}

func extractTokenFromContext(ctx context.Context) string {
	val := ctx.Value(TokenKey)
	if val == nil {
		return ""
	}

	return val.(string)
}

func (h *Handler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	user := &models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Role:      "User",
		Avatar:    req.Avatar,
	}

	createdUser, err := h.service.Register(user)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Uuid:      createdUser.ID,
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Role:      createdUser.Role,
			Avatar:    createdUser.Avatar,
		},
	}, nil
}

// login
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, sessionId, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Uuid:      user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Avatar:    user.Avatar,
		},
		SessionId: sessionId,
	}, nil
}

func (h *Handler) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {

	sessionID := extractSessionFromContext(ctx)
	if sessionID == "" {
		return nil, errors.New("session not found")
	}

	user, err := h.service.GetMe(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &pb.GetMeResponse{
		User: &pb.User{
			Uuid:      user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Avatar:    user.Avatar,
		},
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {

	token := extractTokenFromContext(ctx)

	err := h.service.Logout(ctx, token)
	if err != nil {
		return nil, err
	}

	return &pb.LogoutResponse{}, nil
}

func (h *Handler) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {

	sessionID := extractSessionFromContext(ctx)

	err := h.service.SendEmail(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &pb.SendEmailResponse{
		Success: true,
		Message: "Email sent",
	}, nil
}
