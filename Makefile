up_build: down
	docker-compose up --build -d

up: 
	docker-compose up -d

down:
	docker-compose down

buf:
	docker run --volume $(PWD)/proto:/workspace --workdir /workspace bufbuild/buf generate
