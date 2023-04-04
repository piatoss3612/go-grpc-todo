# go-grpc-todo

## 💻 Buf CLI with Docker

### Initialize Buf module

```powershell
$ docker run --volume "$(pwd)/proto:/workspace" --workdir /workspace bufbuild/buf mod init
```

### Update dependencies

```powershell
$ docker run --volume "$(pwd)/proto:/workspace" --workdir /workspace bufbuild/buf mod update
```

### Generate code

```powershell
$ docker run --volume "$(pwd)/proto:/workspace" --workdir /workspace bufbuild/buf generate
```

```bash
make buf
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