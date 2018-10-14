all: build

TAG?=v0.1.0
IMG?=recommender:${TAG}

build: clean fmt
	go build -o bin/recommender github.com/angao/recommender/cmd/recommender/

clean:
	rm -f bin/recommender

fmt:
	go fmt ./pkg/... ./cmd/...

docker-build:
	docker build . -t ${IMG}

docker-push:
	docker push ${IMG}

.PHONY: all build clean fmt