package messenger

import (
	"context"
	"log"

	"github.com/ros3n/hes/api/models"
	"google.golang.org/grpc"

	pb "github.com/ros3n/hes/lib/communication"
)

type MessageSender interface {
	SendEmail(ctx context.Context, email *models.Email) error
}

type GRPCMessageSender struct {
	address string
}

func NewGRPCMessageSender(addr string) *GRPCMessageSender {
	return &GRPCMessageSender{address: addr}
}

func (gms *GRPCMessageSender) SendEmail(ctx context.Context, email *models.Email) error {
	conn, err := grpc.Dial(gms.address, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()
	client := pb.NewMailerClient(conn)
	_, err = client.SendEmail(ctx, serializedEmail(email))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func serializedEmail(email *models.Email) *pb.SendEmailRequest {
	return &pb.SendEmailRequest{
		Id: email.ID, Sender: email.Sender, Recipients: email.Recipients, Subject: email.Subject, Message: email.Message,
	}
}
