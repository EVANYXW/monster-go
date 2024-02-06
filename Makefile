proto:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go  build -o monsterGo
world: proto
	./monsterGo run --server_name=world