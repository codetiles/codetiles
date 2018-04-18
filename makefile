run:
	echo "Starting codetiles server"
	go run src/*.go

# Linux and MacOS only (currently)
build:
	echo "Building the codetiles server into "
	mkdir -p bin
	$$(cd src && go build -o ../bin/codetiles)
	echo "Executable built in bin/codetiles. Run with ./bin/codetiles"
