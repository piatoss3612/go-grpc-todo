SERVER=server
PROXY=proxy

build: build_server build_proxy

build_server:
	cd cmd/server && set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o ../../build/${SERVER}/${SERVER} ./

build_proxy:
	cd cmd/proxy && set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o ../../build/${PROXY}/${PROXY} ./

server:
	go run ./cmd/server/main.go -p 8080

http:
	go run ./cmd/server/main.go -p 8080 -s http

proxy:
	go run ./cmd/proxy/main.go -p 8081 -e localhost:8080

client:
	go run ./cmd/client/main.go -p 8080

generate:
	docker run --volume $(PWD):/workspace --workdir /workspace bufbuild/buf generate
