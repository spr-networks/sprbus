# sprbus

This is https://github.com/moby/pubsub
with support for unix sockets and grpc

# TODO

use this i go.mod:

```
replace github.com/spr-networks/super/pkg/sprbus v0.0.1 => ../../pkg/sprbus
```

will have to solve build in docker, copy / link libs
but this makes local dev easier

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

**install protoc in ubuntu**

```sh
sudo apt-get install build-essential
wget https://github.com/google/protobuf/releases/download/v2.6.1/protobuf-2.6.1.tar.gz
tar -zxvf protobuf-2.6.1.tar.gz && cd protobuf-2.6.1/
./configure
make -j$(nproc) && make check
make install
protoc --version
```

# Usage

see example/main.go

## command line tool examples

see cmd/main.go

**www logs**
```sh
go run main.go -t www
```

**device and wifi events**
```sh
go run main.go -t device,wifi
```

**network traffic**
```sh
go run main.go -t nft
```

**network traffic in json, no timeout and pipe to jq**
```sh
go run main.go -t nft -j --timeout 0 | jq .
```

# Example code

using default sprbus:
```golang
//publish
sprbus.Publish("wifi:station:event", "{\"json\": \"data\"}")

//subscribe
go sprbus.HandleEvent("wifi", func (topic string, json string) {
    fmt.Println("wifi event", topic, json)
})
```

using a custom unix socket server and client:

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
