server:
	go run ./cmd/server/main.go -p 8080

http:
	go run ./cmd/server/main.go -p 8080 -s http

proxy:
	go run ./cmd/server/main.go -p 8081 -s proxy -e localhost:8080

client:
	go run ./cmd/client/main.go -p 8080

generate:
	docker run --volume $(PWD):/workspace --workdir /workspace bufbuild/buf generate
