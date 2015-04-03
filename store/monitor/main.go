package main

import (
	"os"

	"github.com/deis/deis/store/monitor/bindata"

	"github.com/deis/deis/pkg/boot"
	Log "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 6789
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
	numStores := Getopt("NUM_STORES", "3")
	pgNum := Getopt("PG_NUM", "128")
	keys := make(map[string]string)
	keys[etcdPath+"/size"] = numStores
	keys[etcdPath+"/minSize"] = "1"
	keys[etcdPath+"/pgNum"] = pgNum
	keys[etcdPath+"/delayStart"] = "15"
	return keys
}

func (cb *ControllerBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	setupParams := make(map[string]string)
	setupParams["ETCD_PATH"] = currentBoot.EtcdPath
	setupParams["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort
	setupParams["HOST"] = currentBoot.Host.String()

	createParams := make(map[string]string)
	createParams["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort

	return []*types.Script{
		&types.Script{Name: "bash/setup-monitor.bash", Params: setupParams, Content: bindata.Asset},
		&types.Script{Name: "bash/create-monitor.bash", Params: createParams, Content: bindata.Asset},
	}
}

func (cb *ControllerBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-monitor: starting...")
}

func (cb *ControllerBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	hostname, _ := os.Hostname()
	cmd, args := BuildCommandFromString("/usr/bin/ceph-mon -d -i " + hostname + " --public-addr " + hostname + ":" + string(servicePort))
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (cb *ControllerBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (cb *ControllerBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *ControllerBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-monitor: running...")
}

func (cb *ControllerBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (cb *ControllerBoot) UseConfd() bool {
	return true
}
