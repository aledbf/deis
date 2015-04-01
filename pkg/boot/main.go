//go:generate go-extpoints
package boot

import (
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/deis/deis/pkg/boot/extpoints"

	"github.com/deis/deis/pkg/confd"
	"github.com/deis/deis/pkg/etcd"
	Log "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/os"
	//. "github.com/deis/deis/pkg/net"
	"github.com/deis/deis/pkg/types"
	"github.com/robfig/cron"
	_ "net/http/pprof"
)

const (
	timeout time.Duration = 10 * time.Second
	ttl     time.Duration = timeout * 2
)

var (
	signalChan  = make(chan os.Signal, 2)
	log         = Log.New()
	bootProcess = extpoints.BootComponents
)

// Register register an externsion to be used with this application
func Register(component interface{}, name string) []string {
	return extpoints.Register(component, name)
}

// Start initiate the boot process of the current component
func Start(etcdPath, externalPort string) {
	component, exists := bootProcess.Lookup("deis-component")
	if !exists {
		log.Error("error loading deis extension...")
		log.Infof("%v", bootProcess.Names())
		os.Exit(1)
	}

	log.Info("starting deis component...")

	host := Getopt("HOST", "127.0.0.1")
	etcdPort := Getopt("ETCD_PORT", "4001")
	etcdHostPort := host + ":" + etcdPort
	etcdClient := etcd.NewClient([]string{"http://" + etcdHostPort})

	currentBoot := &types.CurrentBoot{
		EtcdClient: etcdClient,
		EtcdPath:   etcdPath,
		EtcdPort:   etcdPort,
		Host:       net.ParseIP(host),
		Timeout:    timeout,
		TTL:        timeout * 2,
		Port:       externalPort,
	}

	if os.Getenv("DEBUG") != "" {
		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()
	}

	for _, key := range component.MkdirsEtcd() {
		etcd.Mkdir(etcdClient, key)
	}

	for key, value := range component.EtcdDefaults() {
		etcd.SetDefault(etcdClient, key, value)
	}

	component.PreBoot(currentBoot)

	if component.UseConfd() {
		// wait until etcd has discarded potentially stale values
		time.Sleep(timeout + 1)

		// wait for confd to run once and install initial templates
		confd.WaitForInitialConf(signalChan, etcdHostPort, timeout)
	}

	log.Debug("running pre boot scripts")
	preBootScripts := component.PreBootScripts(currentBoot)
	for _, script := range preBootScripts {
		err := RunScript(script.Name, script.Params, script.Content)
		if err != nil {
			log.Printf("command finished with error: %v", err)
			signalChan <- syscall.SIGTERM
		}
	}

	if component.UseConfd() {
		// spawn confd in the background to update services based on etcd changes
		go confd.Launch(signalChan, etcdHostPort)
	}

	log.Debug("running boot daemons")
	servicesToStart := component.BootDaemons(currentBoot)
	for _, daemon := range servicesToStart {
		go RunProcessAsDaemon(signalChan, daemon.Command, daemon.Args)
	}

	portsToWaitFor := component.WaitForPorts()
	log.Debugf("waiting for a service in the port %v", portsToWaitFor)
	for _, portToWait := range portsToWaitFor {
		log.Infof("%v", portToWait)
		//WaitForPort("tcp", "127.0.0.1", string(portToWait), timeout)
	}

	log.Debug("starting periodic publication in etcd...")
	log.Debugf("etcd publication path %s, host %s and port %v", etcdPath, host, externalPort)
	go etcd.PublishService(etcdClient, host, etcdPath, externalPort, uint64(ttl.Seconds()), timeout)

	log.Printf("running post boot scripts")
	postBootScripts := component.PostBootScripts(currentBoot)
	for _, script := range postBootScripts {
		err := RunScript(script.Name, script.Params, script.Content)
		if err != nil {
			log.Printf("command finished with error: %v", err)
			signalChan <- syscall.SIGTERM
		}
	}

	component.PostBoot(currentBoot)

	log.Debug("checking for cron tasks...")
	crons := component.ScheduleTasks(currentBoot)
	for _, cronTask := range crons {
		_cron := cron.New()
		_cron.AddFunc(cronTask.Frequency, cronTask.Code)
		_cron.Start()
	}

	<-signalChan
}
