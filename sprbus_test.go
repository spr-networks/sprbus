package sprbus

import (
	"fmt"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

var sprbusServer *Server

func run_sprbus_server(socket string) {
	//fmt.Printf("test: %v\n", sprbusServer)
	if (sprbusServer != nil) {
		fmt.Println("got server:", sprbusServer)
		return
	}

	sprbusServer, err := NewServer(socket)
	if err != nil {
		log.Fatal(err)
	}
	// does not return

	fmt.Println("server:", sprbusServer)
}

func TestConnect(t *testing.T) {
	socket := "/tmp/spr-test-connect.sock"
	time.Sleep(time.Second)
	go run_sprbus_server(socket)

	// client
	var client *Client
	client, err := NewClient(socket)
	defer client.Close()

	if err != nil {
		t.Fatalf("newClient error: %v", err)
	}

	fmt.Println("cool")
}

func TestPubSub(t *testing.T) {
	socket := "/tmp/spr-test-pubsub.sock"
	go run_sprbus_server(socket)

	// lazy way to wait for server to be up
	time.Sleep(time.Second / 2)
	
	// client
	var client *Client
	client, err := NewClient(socket)
	defer client.Close()

	if err != nil {
		t.Fatalf("newClient error: %v", err)
	}

	var wg sync.WaitGroup

	topicSpr := "spr:group:"

	stream, err := client.SubscribeTopic(topicSpr)
	if nil != err {
		t.Fatalf("Client.SubscribeTopic error: %v", err)
	}

	sendMessages := 0
	gotMessages := 0

	go func() {
		wg.Add(1)
		for {
			reply, err := stream.Recv()
			if io.EOF == err {
				break
			}

			// can get this if client is closing
			if nil != err {
				t.Fatalf("Client recv error: %v", err) // Cancelled desc
			}

			topic := reply.GetTopic()
			value := reply.GetValue()

			//fmt.Printf("sub:reply: %v\n", reply)
			fmt.Printf("sub:topic: %v, sub:value: %v\n", topic, value)

			// verify value is json
			if (value[0] != '{' && value[0] != '[' && value[0] != '"') {
				t.Fatalf("invalid value: %v", value)
			}

			if (len(topic) <= len(topicSpr)) {
				t.Fatalf("invalid topic: %v, subscribe: %v", topic, topicSpr)
			}
			
			gotMessages += 1
			//wg.Done()
		}
	}()

	// lazy way to wait for subscribe to register
	time.Sleep(time.Second / 4)

	// publish some msgs
	for i := 0; i < 5; i++ {
		_, err = client.Publish("spr:group:test", "{\"message\": \"test\"}")
		if err != nil {
			t.Fatalf("publish error: %v", err)
		}

		sendMessages += 1
	}

	for i := 0; i < 5; i++ {
		_, err = client.Publish("spr:group:test", "[1,2,23]")
		if err != nil {
			t.Fatalf("publish error: %v", err)
		}

		sendMessages += 1
	}

	for i := 0; i < 5; i++ {
		_, err = client.Publish("spr:group:test", "\"strstr\"")
		if err != nil {
			t.Fatalf("publish error: %v", err)
		}

		sendMessages += 1
	}

	// send msg we dont subscribe to
	for i := 0; i < 5; i++ {
		_, err = client.Publish("rpc:group:test", "{\"message\": \"test\"}")
		if err != nil {
			t.Fatalf("publish error: %v", err)
		}
	}

	// make sure we have time to receive the msg
	time.Sleep(time.Second / 2)
	
	if (gotMessages != sendMessages) {
		t.Fatalf("invalid num messages received: %v/%v", gotMessages, sendMessages)
	}

	//wg.Wait()
}