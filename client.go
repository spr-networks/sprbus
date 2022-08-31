package sprbus

import (
	"context"
	"log"
	"google.golang.org/grpc"

	pb "github.com/spr-networks/sprbus/pubservice"
)

var ServerEventSock = "/state/api/eventbus.sock"

// Client - object capable of subscribing to a remote event bus
type Client struct {
        path     string
        conn     *grpc.ClientConn
        service  pb.PubsubServiceClient
}

func NewClient(socketPath string) (*Client, error) {
	if socketPath == "" {
		socketPath = ServerEventSock
	}

	conn, err := grpc.Dial("unix:///"+socketPath, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	//defer conn.Close()

        client := new(Client)
        client.path = socketPath
	client.conn = conn
	client.service = pb.NewPubsubServiceClient(conn)

	return client, nil
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) Publish(topic string, value string) (*pb.String, error) {
	return client.service.Publish(context.Background(), &pb.String{Topic: topic, Value: value})
}

func (client *Client) Subscribe(topic string, opts ...grpc.CallOption) (pb.PubsubService_SubscribeTopicClient, error) {
	return client.service.Subscribe(
		context.Background(), &pb.String{Topic: topic},
	)
}

func (client *Client) SubscribeTopic(topic string, opts ...grpc.CallOption) (pb.PubsubService_SubscribeClient, error) {
	return client.service.SubscribeTopic(
		context.Background(), &pb.String{Topic: topic},
	)
}
