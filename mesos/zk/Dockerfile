FROM golang:1.3

WORKDIR /go/src/github.com/deis/deis/mesos/zk

ADD main.go /go/src/github.com/deis/deis/mesos/zk/main.go

RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/deis/deis/mesos/zk

