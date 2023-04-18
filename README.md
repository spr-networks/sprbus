# sprbus

![version](https://img.shields.io/github/v/tag/spr-networks/sprbus?sort=semver&label=version)
[![Go Report Card](https://goreportcard.com/badge/github.com/spr-networks/sprbus)](https://goreportcard.com/report/github.com/spr-networks/sprbus)

This package is a pubsub service using [moby pubsub](https://github.com/moby/pubsub) with support for unix sockets and grpc.

## Command line tools

The client code in `cmd/` can be used to connect to a remote spr api using websockets, or a local unix socket.

![sprbus intro](https://user-images.githubusercontent.com/37542945/231619971-96b18ec8-36a9-4e36-bf37-0b0f1e982c7d.gif)

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

## Custom unix socket server and client

See [example/main.go](https://github.com/spr-networks/sprbus/blob/main/example/main.go) for code to setup a custom unix socket server and client

# TODO

use this i go.mod:

```
replace github.com/spr-networks/super/pkg/sprbus v0.0.1 => ../../pkg/sprbus
```

will have to solve build in docker, copy / link libs
but this makes local dev easier

## cli

- for now publish only works for local connections
- temp enable \* notifications:
  - send request to /notifications with `prefix:"", Notifications: False`
  - disable when we're done to not send excessive data over ws.

### Thank you

The command line tool is built using [https://github.com/charmbracelet/bubbletea](BubbleTea), an awesome TUI Framework.
