# sprbus

![version](https://img.shields.io/github/v/tag/spr-networks/sprbus?sort=semver&label=version)
[![Go Report Card](https://goreportcard.com/badge/github.com/spr-networks/sprbus)](https://goreportcard.com/report/github.com/spr-networks/sprbus)

This package is a pubsub service using [moby pubsub](https://github.com/moby/pubsub) with support for unix sockets and grpc.

## Command line tools

The client code in `cmd/` can be used to connect to a remote spr api using websockets, or a local unix socket.

![sprbus](https://user-images.githubusercontent.com/37542945/232639810-7e17380c-42ea-480b-811e-cf5add04a0d2.gif)

See [cmd/main.go](https://github.com/spr-networks/sprbus/blob/main/cmd/main.go)

```sh
cd cmd/; make
./sprbus --help
```

**remote**

```sh
export TOKEN="SPR-API-TOKEN"
./sprbus --addr 192.168.2.1
```

**local**

```sh
./sprbus
```

**example topics**

```sh
#www and api logs
./sprbus -t log
# device and wifi events
./sprbus -t device,wifi
#network traffic
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

![sprbus intro](https://user-images.githubusercontent.com/37542945/231619971-96b18ec8-36a9-4e36-bf37-0b0f1e982c7d.gif)
example showing how to publish events.

# Development

See [example/main.go](https://github.com/spr-networks/sprbus/blob/main/example/main.go)

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

**Custom unix socket server and client**

See [example/main.go](https://github.com/spr-networks/sprbus/blob/main/example/main.go) for code to setup a custom unix socket server and client

# TODO

**client**
- for now publish only works for local connections

### Thank you

The command line tool is built using [BubbleTea](https://github.com/charmbracelet/bubbletea), an awesome TUI Framework.
