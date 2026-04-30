package grpc

import (
	"context"

	"email_service/internal/services"
	pb "email_service/proto"
)

type Handler struct {
	pb.UnimplementedEmailServiceServer
	service *services.EmailService
}

func NewHandler(s *services.EmailService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {

	err := h.service.SendEmail(ctx, req.To, req.Subject, req.Body)
	if err != nil {
		return &pb.SendEmailResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SendEmailResponse{
		Success: true,
		Message: "Email sent successfully",
	}, nil
}

