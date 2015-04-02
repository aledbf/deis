package main

import (
	"os"

	"github.com/deis/deis/store/metadata/bindata"

	"github.com/deis/deis/pkg/boot"
	Log "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 6800
)

var (
	log          = Log.New()
	etcdPath     = Getopt("ETCD_PATH", "/deis/store")
	externalPort = Getopt("EXTERNAL_PORT", string(servicePort))
)

func init() {
	boot.RegisterComponent(new(ControllerBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type ControllerBoot struct{}

func (cb *ControllerBoot) MkdirsEtcd() []string {
	return []string{etcdPath}
}

func (cb *ControllerBoot) EtcdDefaults() map[string]string {
	return map[string]string{}
}

func (cb *ControllerBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	setupParams := make(map[string]string)
	setupParams["ETCD_PATH"] = currentBoot.EtcdPath
	setupParams["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort
	setupParams["HOST"] = currentBoot.Host.String()
	hostname, _ := os.Hostname()
	setupParams["HOSTNAME"] = hostname

	return []*types.Script{
		&types.Script{Name: "bash/setup-metadata.bash", Params: setupParams, Content: bindata.Asset},
	}
}

func (cb *ControllerBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-metadata: starting...")
}

func (cb *ControllerBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	hostname, _ := os.Hostname()
	cmd, args := BuildCommandFromString("/usr/bin/ceph-mds -d -i " + hostname)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (cb *ControllerBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (cb *ControllerBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *ControllerBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-metadata: running...")
}

func (cb *ControllerBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (cb *ControllerBoot) UseConfd() bool {
	return true
}
