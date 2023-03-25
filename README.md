# go-grpc-todo

## Buf CLI with Docker

### Initialize Buf module

```bash
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf mod init
```

### Lint protobuf

```bash
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf lint
```

### Generate code

```bash
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf generate
```