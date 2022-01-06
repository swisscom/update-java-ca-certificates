build:
	mkdir -p ./bin
	CGO_ENABLED=0 go build -o ./bin/update-java-ca-certificates ./

# Docker Image
# Sorry external user, this is Swisscom-internal only :(
IMAGE=tools-docker-local.artifactory.swisscom.com/swisscom/update-java-ca-certificates
VERSION=0.0.1

docker-build:
	docker build . -t "$(IMAGE):$(VERSION)"

docker-push:
	docker push "$(IMAGE):$(VERSION)"

docker-run:
	docker run --rm "$(IMAGE):$(VERSION)" -h