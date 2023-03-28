# sprbus

This is https://github.com/moby/pubsub
with support for unix sockets and grpc

## update proto

```sh
cd pubservice
protoc -I/usr/local/include -I. \
    -I$GOPATH/pkg/mod \
    -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
    --grpc-gateway_out=. --go_out=plugins=grpc:.\
    --swagger_out=. \
    pubservice.proto
```

# Usage

```go
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
var wg sync.WaitGroup

const (
	sockPath = "/tmp/grpc.sock"
)

func spr_cli() {
	client, err := sprbus.NewClient(sockPath)
	defer client.Close()

	if err != nil {
		log.Fatal("err", err)
	}

	fmt.Println("client connected:", client)

	stream, err := client.SubscribeTopic("spr:test")
	if nil != err {
		log.Fatal(err)
	}

	go func() {
		fmt.Println("recv")
		wg.Add(1)

		i := 0

		for {
			reply, err := stream.Recv()
			if io.EOF == err {
				break
			}

			if nil != err {
				log.Fatal("ERRRRRR ", err) // Cancelled desc
			}

			i += 1

			value := reply.GetValue()

			fmt.Printf("<< sub:reply: %v\n", reply)
			fmt.Println("<< sub:value:", value)
			if i >= 3 {
				wg.Done()
				return
			}
		}
	}()

	for i := 1; i < 5; i++ {
		time.Sleep(time.Second/4)

		fmt.Println(">> pub", "spr:test")
		_, err = client.Publish("spr:test", "samplemsg")
		if err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()

	fmt.Println("done")
}

func spr_server() {
	fmt.Println("server listening...")

	_, err := sprbus.NewServer(sockPath)
	if err != nil {
		log.Fatal(err)
	}

	// does not return
}

func main() {
	fmt.Println("main")

	go spr_server()
	time.Sleep(time.Second/4)
	spr_cli()

}
```

# TODO

change ServerEventSock location default value
