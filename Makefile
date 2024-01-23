.PHONY:
docker:
	docker build -t word-of-wisdom:latest .
	docker build -t word-of-wisdom-client:latest -f Dockerfile.client .

.PHONY:
example: docker
	docker network create -d bridge mynetwork
	docker run -d --name word-of-wisdom-server --net=mynetwork word-of-wisdom:latest --loglevel debug
	docker run --name word-of-wisdom-client --net=mynetwork word-of-wisdom-client:latest --addr word-of-wisdom-server.mynetwork:8080
	docker stop word-of-wisdom-server word-of-wisdom-client
	docker rm word-of-wisdom-server word-of-wisdom-client
	docker network rm mynetwork

.PHONY:
lint:
	golangci-lint run --fix
	go test -race ./...
