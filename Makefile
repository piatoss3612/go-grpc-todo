server:
	go run ./cmd/server/main.go -p 8080

http:
	go run ./cmd/server/main.go -p 8080 -s http

client:
	go run ./cmd/client/main.go -p 8080

generate:
	docker run --volume $(PWD):/workspace --workdir /workspace bufbuild/buf generate
