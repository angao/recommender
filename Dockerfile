# Build the manager binary
FROM golang:1.10.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/angao/recommender
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/
COPY version/ version/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o recommender github.com/angao/recommender/cmd/recommender

FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /go/src/github.com/angao/recommender .
CMD ./recommender --v=4 --stderrthreshold=info --db-config-file="/etc/db.yaml" --prometheus-address=http://192.168.99.100:30003
