package messenger

import (
	"context"
	"log"
	"net"

	pb "github.com/ros3n/hes/lib/communication"
	"github.com/ros3n/hes/mailer/models"
	"google.golang.org/grpc"
)

type MessageReceiver interface {
	Start(chan<- *models.Email) error
	Stop()
}

type GRPCMessageReceiver struct {
	address string
	server  *grpc.Server
}

func NewGRPCMessageReceiver(address string) *GRPCMessageReceiver {
	return &GRPCMessageReceiver{address: address}
}

func (gm *GRPCMessageReceiver) Start(newEmailsChan chan<- *models.Email) error {
	lis, err := net.Listen("tcp", gm.address)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return err
	}

	gm.server = grpc.NewServer()
	pb.RegisterMailerServer(gm.server, &mailerServer{newEmailsChan: newEmailsChan})

	go func() {
		if err := gm.server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (gm *GRPCMessageReceiver) Stop() {
	gm.server.GracefulStop()
}

type mailerServer struct {
	newEmailsChan chan<- *models.Email
}

func (ms *mailerServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailReply, error) {
	email := parseSendEmailRequest(req)
	ms.newEmailsChan <- email
	return &pb.SendEmailReply{}, nil
}

func parseSendEmailRequest(req *pb.SendEmailRequest) *models.Email {
	return &models.Email{
		ID:         req.GetId(),
		Sender:     req.GetSender(),
		Recipients: req.GetRecipients(),
		Subject:    req.GetSubject(),
		Message:    req.GetMessage(),
	}
}
