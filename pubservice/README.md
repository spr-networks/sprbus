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
