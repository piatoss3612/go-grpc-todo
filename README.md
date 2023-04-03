# go-grpc-todo

## 💻 Buf CLI with Docker

### Initialize Buf module

```powershell
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf mod init
```

### Generate code

```powershell
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf generate
```

### Update dependencies

```powershell
$ docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf mod update
```

## 🐳 Docker Compose

### Build and run

```
make up_build
```

### Open browser

```
http://localhost:80
```