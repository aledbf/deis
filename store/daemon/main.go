package main

import (
	"os"

	"github.com/deis/deis/store/daemon/bindata"

	"github.com/deis/deis/pkg/boot"
	"github.com/deis/deis/pkg/etcd"
	Log "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = -1
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
		&types.Script{Name: "bash/setup-daemon.bash", Params: setupParams, Content: bindata.Asset},
	}
}

func (cb *ControllerBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-daemon: starting...")
}

func (cb *ControllerBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	osdID := etcd.Get(currentBoot.EtcdClient, "/deis/store/osds/"+currentBoot.Host.String())
	cmd, args := BuildCommandFromString("ceph-osd -d -i " + osdID + " -k /var/lib/ceph/osd/ceph-" + osdID + "/keyring")
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (cb *ControllerBoot) WaitForPorts() []int {
	return []int{}
}

func (cb *ControllerBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *ControllerBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-daemon: radosgw running...")
}

func (cb *ControllerBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (cb *ControllerBoot) UseConfd() bool {
	return true
}
