package main

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"
	"github.com/spr-networks/sprbus"
)

var client *sprbus.Client
//var socket = "/tmp/test.sock"
var socket = "/state/api/eventbus.sock"

func do_subscribe(client *sprbus.Client) {
	time.Sleep(time.Second / 4)

	var wg sync.WaitGroup

	stream, err := client.SubscribeTopic("spr:test")
	if nil != err {
		log.Fatal(err)
	}

	go func() {
		wg.Add(1)
		for {
			reply, err := stream.Recv()
			if io.EOF == err {
				break
			}

			if nil != err {
				return
			}

			topic := reply.GetTopic()
			value := reply.GetValue()

			fmt.Printf("topic=%v value=%v\n", topic, value)
		}
	}()

}

func do_publish(client *sprbus.Client) {
	for i := 0; i < 5; i++ {
		/*_, err := client.Publish("spr:test", "{\"data\": \"test\"}")
		if err != nil {
			log.Fatal(err)
		}*/

		type testS struct {
			Title string
			Body string
		}

		sprbus.Publish("spr:test:struct", testS{Title: "tttt1111", Body: "datahere"})
		sprbus.Publish("spr:test:string", "s1ACID")
		sprbus.Publish("spr:test:array", []int{11,23})
	}
}

func spr_server() {
	fmt.Println("server listening...")

	server, err := sprbus.NewServer(socket)
	if err != nil {
		log.Fatal(err)
	}

	// does not return
	fmt.Println("server:", server)
}

func main() {
	go spr_server()

	time.Sleep(time.Second / 4)

	client, err := sprbus.NewClient(socket)
	defer client.Close()

	if err != nil {
		log.Fatal("err", err)
	}

	fmt.Println("client:", client)

	do_subscribe(client)

	do_publish(client)

	fmt.Println("done")
}
