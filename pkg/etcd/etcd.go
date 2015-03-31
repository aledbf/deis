package etcd

import (
	"errors"
	"time"

	"github.com/coreos/go-etcd/etcd"
	Log "github.com/deis/deis/pkg/log"
)

var log = Log.New()

func NewClient(machines []string) *etcd.Client {
	return etcd.NewClient(machines)
}

func SetDefault(client *etcd.Client, key, value string) {
	Set(client, key, value, 0)
}

func Mkdir(client *etcd.Client, path string) {
	_, err := client.CreateDir(path, 0)
	if err != nil {
		log.Debug(err)
	}
}

// WaitForKeysEtcd wait for the required keys up to the timeout or forever if is nil
func WaitForKeys(client *etcd.Client, keys []string, ttl time.Duration) error {
	start := time.Now()
	wait := true

	for {
		for _, key := range keys {
			_, err := client.Get(key, false, false)
			if err != nil {
				log.Debugf("key \"%s\" error %v", key, err)
				wait = true
			}
		}

		if !wait {
			return nil
		}

		log.Debug("waiting for missing etcd keys...")
		time.Sleep(1 * time.Second)
		wait = false

		if time.Since(start) > ttl {
			return errors.New("maximum ttl reached. aborting")
		}
	}
}

func Get(client *etcd.Client, key string) string {
	result, err := client.Get(key, false, false)
	if err != nil {
		log.Debugf("%v", err)
		return ""
	}

	return result.Node.Value
}

func GetList(client *etcd.Client, key string) []string {
	values, err := client.Get(key, true, false)
	if err != nil {
		log.Debugf("%v", err)
		return []string{}
	}

	result := []string{}
	for _, node := range values.Node.Nodes {
		result = append(result, node.Value)
	}

	log.Infof("%v", result)
	return result
}

func Set(client *etcd.Client, key, value string, ttl uint64) {
	_, err := client.Set(key, value, ttl)
	if err != nil {
		log.Debugf("%v", err)
	}
}

// Publish a service to etcd periodically
func PublishService(
	client *etcd.Client,
	host string,
	etcdPath string,
	externalPort string,
	ttl uint64,
	timeout time.Duration) {

	for {
		Set(client, etcdPath+"/host", host, ttl)
		Set(client, etcdPath+"/port", externalPort, ttl)
		time.Sleep(timeout)
	}
}
