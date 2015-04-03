package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/deis/deis/builder/bindata"

	"github.com/deis/deis/pkg/boot"
	Log "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 22
)

var (
	etcdPath     = Getopt("ETCD_PATH", "/deis/builder")
	externalPort = Getopt("EXTERNAL_PORT", string(servicePort))
	log          = Log.New()
)

func init() {
	boot.RegisterComponent(new(BuilderBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type BuilderBoot struct{}

func (bb *BuilderBoot) MkdirsEtcd() []string {
	return []string{
		etcdPath,
		etcdPath + "/users",
	}
}

func (bb *BuilderBoot) EtcdDefaults() map[string]string {
	return map[string]string{}
}

func (bb *BuilderBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{
		&types.Script{Name: "bash/check-overlay.bash", Content: bindata.Asset},
	}
}

func (bb *BuilderBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-builder: starting...")
	// remove any pre-existing docker.sock
	os.Remove("/var/run/docker.sock")
}

func (bb *BuilderBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	driverOverride := readDockerEnvFile()
	log.Debugf("custom docker env [%v]", driverOverride)
	docker := "docker -D -d --bip=172.19.42.1/16 " +
		driverOverride +
		" --insecure-registry 10.0.0.0/8 " +
		" --insecure-registry 172.16.0.0/12 " +
		" --insecure-registry 192.168.0.0/16 " +
		" --insecure-registry 100.64.0.0/10"

	log.Debugf("starting docker daemon: %v", docker)
	dockerCmd, dockerArgs := BuildCommandFromString(docker)
	sshCmd, sshArgs := BuildCommandFromString("/usr/sbin/sshd -D -e")
	return []*types.ServiceDaemon{
		&types.ServiceDaemon{Command: dockerCmd, Args: dockerArgs},
		&types.ServiceDaemon{Command: sshCmd, Args: sshArgs},
	}
}

func (bb *BuilderBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (bb *BuilderBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (bb *BuilderBoot) PostBoot(currentBoot *types.CurrentBoot) {
	waitForDocker()
	RunScript("bash/build-slugbuilder-slugrunner.bash", nil, bindata.Asset)

	log.Info("deis-builder: running...")
}

func (bb *BuilderBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (bb *BuilderBoot) UseConfd() bool {
	return true
}

func waitForDocker() {
	log.Debug("waiting for docker daemon to be available...")
	for {
		cmd := exec.Command("docker", "info")
		if err := cmd.Run(); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func readDockerEnvFile() string {
	envFile, err := ioutil.ReadFile("/etc/docker.env")
	if err != nil {
		log.Debugf("%v", err)
		return ""
	}
	return string(envFile)
}
