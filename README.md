mq service
==========
A service which encapsulates all logic for interfacing with our message broker, implemented in Go.

### development
##### protobuf
To generate the needed protobuf code for this service, read the [leaf-proto README](https://gitlab.com/project-leaf/leaf-proto) first. Then do the following:
```bash
protoc -I leaf-proto/ leaf-proto/core.proto --go_out=plugins=grpc:src/proto/core
protoc -I leaf-proto/ leaf-proto/mq-service.proto --go_out=Mcore.proto=gitlab.com/project-leaf/mq-service-go/src/proto/core,plugins=grpc:src/proto/mq
```
