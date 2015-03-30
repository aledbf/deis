package boot

import (
	"net"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

type Boot struct {
	Etcd     *etcd.Client
	EtcdPath string
	EtcdPort string
	Host     net.IP
	Port     string
	Timeout  time.Duration
	TTL      time.Duration
}
