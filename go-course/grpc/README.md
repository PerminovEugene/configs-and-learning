# Golang GRPC Server and Client

## Local setup

1. Install golang, protoc and google.golang.org/grpc libraries

2. Run server

3. Run client

## Regenerate greet.pb.go

```console
   protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
```
