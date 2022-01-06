build:
	mkdir -p ./bin
	CGO_ENABLED=0 go build -o ./bin/update-java-ca-certificates ./