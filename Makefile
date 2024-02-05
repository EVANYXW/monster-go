build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go  build -o monsterGo
world: build
	./monsterGo run --server_name=world