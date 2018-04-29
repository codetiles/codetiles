run: depends
	@echo "Starting codetiles server..."
	@go run src/*.go
build: depends
	@echo "Building the codetiles server into "
	@mkdir -p bin
	$$(cd src && go build -o ../bin/codetiles)
	@echo "Executable built in bin/codetiles. Run with ./bin/codetiles"
windows: depends
	@echo Starting codetiles server...
	@go run src/main.go src/map.go src/game.go src/users.go src/lobby.go src/code.go src/misc.go src/tick.go
depends:
	@go get github.com/gorilla/websocket
