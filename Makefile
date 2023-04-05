up_build: down
	docker-compose up --build -d

up: 
	docker-compose up -d

down:
	docker-compose down

buf:
	docker run --volume $(PWD)/proto:/workspace --workdir /workspace bufbuild/buf generate

dbgen:
	docker run --rm -v $(PWD)/db:/src -w /src kjconroy/sqlc generate