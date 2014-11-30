FROM golang:1.3

WORKDIR /go/src/github.com/deis/deis/mesos/master

ADD main.go /go/src/github.com/deis/deis/mesos/master/main.go

RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/deis/deis/mesos/master

