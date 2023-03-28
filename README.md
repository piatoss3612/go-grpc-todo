# go-grpc-todo

## Buf CLI with Docker

### Initialize Buf module

```powershell
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf mod init
```

### Lint protobuf

```powershell
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf lint
```

### Generate code

```powershell
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf generate
```

## Generate Reverse Proxy

### Update dependencies

```powershell
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf mod update
```

### Add annotations

```proto
import "google/api/annotations.proto";

option go_package = "./todo/v1;todo";

service TodoService {
    rpc Add(AddRequest) returns (AddResponse) {
        option (google.api.http) = {
            post: "/v1/todo"
            body: "*"
        };
    };
    rpc Get(GetRequest) returns (Todo) {
        option (google.api.http) = {
            get: "/v1/todo/{id}"
        };
    };

    ...
}
```

### Add gateway plugin

```yaml
version: v1
plugins:
  - plugin: buf.build/grpc-ecosystem/gateway
    out: gen/go
    opt:
      - paths=source_relative
      - generate_unbound_methods=true
```

### Generate code

```bash
make generate
```

## Proxy Server

### Start GRPC Server

```
make server
```

### Start Proxy Server

```
make proxy
```