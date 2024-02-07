proto:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go  build -o monsterGo
world: proto
	./monsterGo run --server_name=world --env fat
compose-up: Dockerfile
	docker build -t monster-go .
docker-run: compose-up
	docker run -itd --name monster-go -p 8023:8023 -p 8024:8024 -p 6060:6060 monster-go