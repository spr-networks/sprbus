# sprbus

![version](https://img.shields.io/github/v/tag/spr-networks/sprbus?sort=semver&label=version)


This is [moby pubsub](https://github.com/moby/pubsub) with support for unix sockets and grpc

# Usage

see example/main.go

## command line tools

see cmd/main.go

```sh
cd cmd/; make
./sprbus --help
```

**www logs**
```sh
./sprbus -t www
```

**device and wifi events**
```sh
./sprbus -t device,wifi
```

**network traffic**
```sh
./sprbus -t nft
```

**network traffic in json, no timeout and pipe to jq**
```sh
./sprbus -t nft -j --timeout 0 | jq .
```

**publish test event**
```sh
./sprbus -t test:event -p '{"msg": "testevent1234"}'
```

# Example code

using default sprbus:
```golang
//publish json string
sprbus.PublishString("wifi:station:event", "{\"json\": \"data\"}")

//publish object
sprbus.Publish("www:auth:user:fail", map[string]string{"username": username})

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

# Update pubservice proto

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

**build**

```sh
cd pubservice
protoc -I/usr/local/include -I. \
    -I$GOPATH/pkg/mod \
    -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
    --grpc-gateway_out=. --go_out=plugins=grpc:.\
    --swagger_out=. \
    pubservice.proto
```

# TODO

use this i go.mod:

```
replace github.com/spr-networks/super/pkg/sprbus v0.0.1 => ../../pkg/sprbus
```

will have to solve build in docker, copy / link libs
but this makes local dev easier

