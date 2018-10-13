all: build

TAG?=v0.1.0

build: clean fmt
	go build -o bin/recommender github.com/angao/recommender/cmd/recommender/

clean:
	rm -f bin/recommender

fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: all build clean fmt