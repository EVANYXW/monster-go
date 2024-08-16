build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go  build -o monsterGo
world: build
	./monsterGo run --server_name=world --env fat
docker-compose: Dockerfile
	docker build -f Dockerfile -t monster-go .
docker-run: docker-compose
	docker run -itd --name monster-go -p 8023:8023 -p 8024:8024 -p 6060:6060  --restart=always monster-go