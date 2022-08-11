# sprbus

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

# TODO

change ServerEventSock location default value
