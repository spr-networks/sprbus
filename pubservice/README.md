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
protoc -I . \
  --go_out . --go_opt paths=source_relative \
  --go-grpc_out . --go-grpc_opt paths=source_relative \
  --grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
  --openapiv2_out . \
  *.proto
```
