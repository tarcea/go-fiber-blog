dev:
	CompileDaemon -command="./go-fiber-blog"
test:
	go test -v --cover ./...