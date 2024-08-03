package sprbus

import (
	"context"
	"encoding/json"
	pb "github.com/spr-networks/sprbus/pubservice"
	"google.golang.org/grpc"
	"io"
)

var ServerEventSock = "/state/api/eventbus.sock"

// Client - object capable of subscribing to a remote event bus
type Client struct {
	path    string
	conn    *grpc.ClientConn
	service pb.PubsubServiceClient
}

var gClient *Client

func getClient() (*Client, error) {
	if gClient.path == "" {
		client, err := NewClient(ServerEventSock)
		if err != nil {
			return nil, err
		}
		gClient = client
	}
	return gClient, nil
}

// sprbus.PublishString() using default socket
func PublishString(topic string, value string) (*pb.String, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.Publish(topic, value)
}

// sprbus.Publish() using default socket, make sure bytes are json
func Publish(topic string, bytes interface{}) (*pb.String, error) {
	value, err := json.Marshal(bytes)

	if err != nil {
		return nil, err
	}

	return PublishString(topic, string(value))
}

func HandleEvent(topic string, callback func(string, string)) error {
	client, err := getClient()

	if nil != err {
		return err
	}

	stream, err := client.SubscribeTopic(topic)
	if nil != err {
		return err
	}

	for {
		reply, err := stream.Recv()
		if io.EOF == err {
			break
		}

		if nil != err {
			// Cancelled desc
			return nil
		}

		topic := reply.GetTopic()
		json := reply.GetValue()

		callback(topic, json)

	}

	return nil
}

func NewClient(socketPath string) (*Client, error) {
	if socketPath == "" {
		socketPath = ServerEventSock
	}

	conn, err := grpc.Dial("unix:///"+socketPath, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

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
