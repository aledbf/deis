package main

import (
	"fmt"
	"os"

	"github.com/deis/deis/store/gateway/bindata"

	"github.com/deis/deis/pkg/boot"
	Log "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 8888
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
		&types.Script{Name: "bash/setup-gateway.bash", Params: setupParams, Content: bindata.Asset},
	}
}

func (cb *ControllerBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-gateway: starting...")
}

func (cb *ControllerBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	cmd, args := BuildCommandFromString("/etc/init.d/radosgw start")
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (cb *ControllerBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (cb *ControllerBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *ControllerBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-gateway: radosgw running...")
}

func (cb *ControllerBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	params := make(map[string]string)
	params["HOSTNAME"] = currentBoot.Host.String()
	params["ETCD_PATH"] = currentBoot.EtcdPath
	params["ETCD_TTL"] = fmt.Sprintf("%v", currentBoot.TTL.Seconds())
	params["EXTERNAL_PORT"] = currentBoot.Port
	params["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort

	return []*types.Cron{
		&types.Cron{
			Frequency: "@every 5s",
			Code: func() {
				RunScript("bash/gateway-master.bash", params, bindata.Asset)
			},
		},
	}
}

func (cb *ControllerBoot) UseConfd() bool {
	return true
}
