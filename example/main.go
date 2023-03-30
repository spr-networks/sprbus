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
var socket = "/tmp/test.sock"

func spr_cli() {

	time.Sleep(time.Second / 4)

	client, err := sprbus.NewClient(socket)
	defer client.Close()

	if err != nil {
		log.Fatal("err", err)
	}

	var wg sync.WaitGroup

	fmt.Println("got client:", client)

	stream, err := client.SubscribeTopic("spr:test")
	if nil != err {
		log.Fatal(err)
	}

	go func() {
		fmt.Println("recv")
		wg.Add(1)
		for {
			reply, err := stream.Recv()
			if io.EOF == err {
				break
			}

			if nil != err {
				//log.Fatal("ERRRRRR ", err) // Cancelled desc
				return
			}

			value := reply.GetValue()

			fmt.Printf("sub:reply: %v\n", reply)
			fmt.Println("sub:value:", value)
			//wg.Done()
			//return
		}
	}()

	for i := 1; i < 2; i++ {

		time.Sleep(time.Second / 4)

		fmt.Println("pub")
		_, err = client.Publish("spr:test", "samplemsg")
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Publish("rpc:test", "samplemsg")
		if err != nil {
			log.Fatal(err)
		}

	}

	//wg.Wait()

	fmt.Println("done")

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
	fmt.Println("main")

	go spr_server()
	spr_cli()

	/*for i := 1; i < 10; i++ {
		time.Sleep(time.Second)
	}*/
}
