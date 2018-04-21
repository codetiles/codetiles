run:
	echo "Starting codetiles server"
	go run src/*.go
build:
	echo "Building the codetiles server into "
	mkdir -p bin
	$$(cd src && go build -o ../bin/codetiles)
	echo "Executable built in bin/codetiles. Run with ./bin/codetiles"
windows:
	echo "Starting codetiles server"
	go run src/main.go src/map.go src/game.go src/users.go
