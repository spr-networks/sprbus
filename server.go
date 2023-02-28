package sprbus

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"github.com/moby/moby/pkg/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/spr-networks/sprbus/pubservice"
)

//defined in client.go
//var ServerEventSock = "/tmp/grpc.sock"

type Server struct {
        path     string
        server   *grpc.Server
}

type PubsubService struct {
	pub *pubsub.Publisher
}

func (p *PubsubService) Publish(ctx context.Context, arg *pb.String) (*pb.String, error) {
	msg := arg.GetTopic() + ":" + arg.GetValue()
	//p.pub.Publish(arg.GetValue())
	p.pub.Publish(msg)
	return &pb.String{}, nil
}

func (p *PubsubService) SubscribeTopic(arg *pb.String, stream pb.PubsubService_SubscribeTopicServer) error {
	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, arg.GetTopic()) {
				return true
			}
		}
		return false
	})

	for v := range ch {
		topic := arg.GetTopic()
		value := strings.TrimPrefix(v.(string), topic+":")
		if err := stream.Send(&pb.String{Topic: topic, Value: value}); nil != err {
			return err
		}
	}

	return nil
}

func (p *PubsubService) Subscribe(arg *pb.String, stream pb.PubsubService_SubscribeServer) error {
	ch := p.pub.Subscribe()
	for v := range ch {
		if err := stream.Send(&pb.String{Value: v.(string)}); nil != err {
			return err
		}
	}
	return nil
}

func NewPubsubService() *PubsubService {
	return &PubsubService{pub: pubsub.NewPublisher(100*time.Millisecond, 10)}
}

func NewServer(socketPath string) (*Server, error) {

	if socketPath == "" {
		socketPath = ServerEventSock
	}

	os.Remove(socketPath)

	lis, err := net.Listen("unix", socketPath)

	if err != nil {
		return nil, err
	}

	server := new(Server)
	server.path = socketPath
	server.server = grpc.NewServer()

	//  register  grpcurl  The required  reflection  service 
	reflection.Register(server.server)

	//  Register pubsub
	pb.RegisterPubsubServiceServer(server.server, NewPubsubService())

	//fmt.Println("starting grpc server...")

	if err := server.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return server, nil
}
