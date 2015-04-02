package main

import (
	"fmt"

	"github.com/deis/deis/registry/bindata"

	"github.com/deis/deis/pkg/boot"
	Log "github.com/deis/deis/pkg/log"
	"github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 5000
)

var (
	etcdPath     = os.Getopt("ETCD_PATH", "/deis/registry")
	externalPort = os.Getopt("EXTERNAL_PORT", string(servicePort))
	log          = Log.New()
)

func init() {
	boot.RegisterComponent(new(RegistryBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type RegistryBoot struct{}

func (rb *RegistryBoot) MkdirsEtcd() []string {
	return []string{
		etcdPath,
		"/deis/store/gateway",
	}
}

func (rb *RegistryBoot) EtcdDefaults() map[string]string {
	bucketName := os.Getopt("BUCKET_NAME", "registry")
	keys := make(map[string]string)
	keys[etcdPath+"/protocol"] = "http"
	keys[etcdPath+"/bucketName"] = bucketName
	return keys
}

func (rb *RegistryBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	bucketName := os.Getopt("BUCKET_NAME", "registry")
	params := make(map[string]string)
	params["BUCKET_NAME"] = bucketName
	return []*types.Script{
		&types.Script{Name: "bash/create-bucket.bash", Params: params, Content: bindata.Asset},
	}
}

func (rb *RegistryBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-registry: starting...")
}

func (rb *RegistryBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	cmd, args := os.BuildCommandFromString("sudo -E -u registry docker-registry")
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (rb *RegistryBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (rb *RegistryBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (rb *RegistryBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-registry: docker-registry is running...")
}

func (rb *RegistryBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	params := make(map[string]string)
	params["HOSTNAME"] = os.Getopt("HOSTNAME", "localhost")
	params["HOST"] = currentBoot.Host.String()
	params["ETCD_PATH"] = currentBoot.EtcdPath
	params["ETCD_TTL"] = fmt.Sprintf("%v", currentBoot.TTL.Seconds())
	params["EXTERNAL_PORT"] = currentBoot.Port
	params["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort

	return []*types.Cron{
		&types.Cron{
			Frequency: "@every 5s",
			Code: func() {
				os.RunScript("bash/registry-master.bash", params, bindata.Asset)
			},
		},
	}
}

func (rb *RegistryBoot) UseConfd() bool {
	return true
}
