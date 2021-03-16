# GRPC

Commands:

```
go get package
```

Install packages on the project

```
go install package
```

Sincronize modules

```
go mod tidy
```

Compile .proto files and generates go stubs

```
protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
```

###Links

[Go lang](https://golang.org/)

[GRPC](https://grpc.io/)

[GRPC Client for tests](https://github.com/ktr0731/evans)

[Protocol Buffers](https://developers.google.com/protocol-buffers)
