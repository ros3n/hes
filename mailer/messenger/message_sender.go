package messenger

import (
	"context"
	"log"

	"github.com/ros3n/hes/mailer/models"
	"google.golang.org/grpc"

	pb "github.com/ros3n/hes/lib/communication"
)

type MessageSender interface {
	SendStatus(ctx context.Context, status *models.SendStatus) error
}

type GRPCMessageSender struct {
	address string
}

func NewGRPCMessageSender(addr string) *GRPCMessageSender {
	return &GRPCMessageSender{address: addr}
}

func (gms *GRPCMessageSender) SendStatus(ctx context.Context, status *models.SendStatus) error {
	conn, err := grpc.Dial(gms.address, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()
	client := pb.NewMailerAPIClient(conn)
	_, err = client.SendStatus(ctx, serializedStatus(status))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func serializedStatus(status *models.SendStatus) *pb.SendStatusRequest {
	return &pb.SendStatusRequest{
		Id: status.EmailID, Success: status.Success,
	}
}
