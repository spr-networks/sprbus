package sprbus

import (
	"fmt"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

func run_sprbus_server(socket string) {
	server, err := NewServer(socket)
	if err != nil {
		log.Fatal(err)
	}
	// does not return

	fmt.Println("server:", server)
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
				//t.Fatalf("Client recv error: %v", err) // Cancelled desc
				return
			}

			topic := reply.GetTopic()
			value := reply.GetValue()

			//fmt.Printf("sub:reply: %v\n", reply)
			fmt.Printf("sub:topic: %v, sub:value: %v\n", topic, value)

			// verify value is json
			if value[0] != '{' && value[0] != '[' && value[0] != '"' {
				t.Fatalf("invalid value: %v", value)
			}

			if len(topic) <= len(topicSpr) {
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

	if gotMessages != sendMessages {
		t.Fatalf("invalid num messages received: %v/%v", gotMessages, sendMessages)
	}

	//wg.Wait()
}

/*func TestHandleEvent(t *testing.T) {
	fmt.Println("need /state/api/ dir")
	go HandleEvent("", func(topic string, json string) {
		fmt.Printf("[sprbus] %v %v\n", topic, json)
	})

	time.Sleep(time.Second/2)
	Publish("spr:test:string", "s1ACID")
	Publish("spr:test:string", "s1ACID")
	Publish("spr:test:string", "s1ACID")
	Publish("spr:test:string", "s1ACID")
	time.Sleep(time.Second/2)
}*/

func TestVerifyTopicWildcard(t *testing.T) {
	socket := "/tmp/spr-test-subw.sock"
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

	//subscribe to everything
	stream, err := client.SubscribeTopic("")
	if nil != err {
		t.Fatalf("Client.SubscribeTopic error: %v", err)
	}

	sendMessages := 0
	gotMessages := 0

	topicTest := "spr:test:wildcard:subscribe"

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

			//verify topic is not empty
			if topic != topicTest {
				t.Fatalf("invalid topic set subscribe: %v, should be: %v", topic, topicTest)
			}

			// verify value is json
			if value[0] != '{' && value[0] != '[' && value[0] != '"' {
				t.Fatalf("invalid value: %v", value)
			}

			gotMessages += 1
			//wg.Done()
		}
	}()

	// lazy way to wait for subscribe to register
	time.Sleep(time.Second / 4)

	// publish some msgs
	for i := 0; i < 5; i++ {
		_, err = client.Publish(topicTest, "{\"message\": \"test\"}")
		if err != nil {
			t.Fatalf("publish error: %v", err)
		}

		sendMessages += 1
	}

	// make sure we have time to receive the msg
	time.Sleep(time.Second / 2)

	if gotMessages != sendMessages {
		t.Fatalf("invalid num messages received: %v/%v", gotMessages, sendMessages)
	}

	//wg.Wait()
}
