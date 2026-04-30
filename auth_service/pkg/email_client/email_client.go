package emailclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"



	pb "email_service/proto"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.EmailServiceClient
}

func NewEmailClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: pb.NewEmailServiceClient(conn),
	}, nil
}

func (c *Client) SendEmail(ctx context.Context, to, subject, body string) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.client.SendEmail(ctx, &pb.SendEmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
	})

	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
