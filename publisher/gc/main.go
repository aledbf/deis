package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

const (
	defaultRefreshTime time.Duration = 40 * time.Second
	defaultEtcdHost                  = "127.0.0.1"
	defaultEtcdPort                  = "4001"
)

var (
	etcdHost = flag.String("etcd-host", defaultEtcdHost, "The etcd host.")
	etcdPort = flag.String("etcd-port", defaultEtcdPort, "The etcd port.")
)

func main() {
	flag.Parse()

	etcdClient := etcd.NewClient([]string{"http://" + *etcdHost + ":" + *etcdPort})

	for {
		go periodicGC(etcdClient)
		time.Sleep(defaultRefreshTime)
	}
}

// periodicGC checks that the published containers are alive
// (checking if the published port is accepting connections)
func periodicGC(client *etcd.Client) {
	data, err := client.Get("/deis/services", true, true)
	if err != nil {
		return
	}
	for _, node := range data.Node.Nodes {
		if node.Dir {
			for _, appNode := range node.Nodes {
				appInstance := appNode.Key
				appHostPort := appNode.Value
				if !isPortOpen(appHostPort) {
					log.Printf("not running '%s', removing\n", appInstance)
					client.Delete(appInstance, false)
				}
			}
		}
	}
}

// isPortOpen checks if the given port is accepting tcp connections
func isPortOpen(hostAndPort string) bool {
	portOpen := false
	conn, err := net.Dial("tcp", hostAndPort)
	if err == nil {
		portOpen = true
		defer conn.Close()
	}
	return portOpen
}
