package messenger

import (
	"context"
	"github.com/ros3n/hes/api/models"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/ros3n/hes/lib/communication"
)

type MessageReceiver interface {
	Start(chan<- *models.SendStatus) error
	Stop()
}

type GRPCMessageReceiver struct {
	address string
	server  *grpc.Server
}

func NewGRPCMessageReceiver(address string) *GRPCMessageReceiver {
	return &GRPCMessageReceiver{address: address}
}

func (gm *GRPCMessageReceiver) Start(sendStatusChan chan<- *models.SendStatus) error {
	lis, err := net.Listen("tcp", gm.address)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return err
	}

	gm.server = grpc.NewServer()
	pb.RegisterMailerAPIServer(gm.server, &mailerAPIServer{sendStatusChan: sendStatusChan})

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

type mailerAPIServer struct {
	sendStatusChan chan<- *models.SendStatus
}

func (ms *mailerAPIServer) SendStatus(ctx context.Context, req *pb.SendStatusRequest) (*pb.SendStatusReply, error) {
	status := parseSendStatusRequest(req)
	ms.sendStatusChan <- status
	return &pb.SendStatusReply{}, nil
}

func parseSendStatusRequest(req *pb.SendStatusRequest) *models.SendStatus {
	return &models.SendStatus{
		ID:      req.GetId(),
		Success: req.GetSuccess(),
	}
}
