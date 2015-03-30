package boot

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deis/deis/pkg/confd"
	"github.com/deis/deis/pkg/etcd"
	. "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/net"
	. "github.com/deis/deis/pkg/os"
)

const (
	timeout time.Duration = 10 * time.Second
	ttl     time.Duration = timeout * 2
	wait    time.Duration = timeout / 2
)

var (
	signalChan = make(chan os.Signal, 2)
)

// New contructor that indicates the etcd base path and
// the port that the component will expose
func New(etcdPath, port string) *Boot {
	Log.Info("starting deis component...")

	host := Getopt("HOST", "127.0.0.1")
	etcdPort := Getopt("ETCD_PORT", "4001")

	etcdHostPort := host + ":" + etcdPort

	etcdClient := etcd.NewClient([]string{"http://" + etcdHostPort})

	signalChan = make(chan os.Signal, 2)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	return &Boot{
		Etcd:     etcdClient,
		EtcdPath: etcdPath,
		EtcdPort: etcdPort,
		Host:     net.ParseIP(host),
		Timeout:  timeout,
		TTL:      timeout * 2,
		Port:     port,
	}
}

// Start initiates the boot process waiting for the correct initialization
// of the required values for the confd template and launch confd as daemon
func (this *Boot) Start() {
	// wait until etcd has discarded potentially stale values
	time.Sleep(this.Timeout + 1)

	// wait for confd to run once and install initial templates
	confd.WaitForInitialConf(signalChan, this.Host.String()+":"+this.EtcdPort, this.Timeout)

	// spawn confd in the background to update services based on etcd changes
	go confd.Launch(signalChan, this.Host.String()+":"+this.EtcdPort)
}

// Publish publish information about the relevant process running in the boot
// process in etcd using specified path and port/s
func (this *Boot) Publish(port ...string) {
	portToPublish := this.Port
	// If we specify a custom port we use that one
	if len(port) != 0 {
		portToPublish = port[1]
	}
	Log.Info("starting periodic publication in etcd...")
	Log.Debugf("etcd publication path %s, host %s and port %s", this.EtcdPath, this.Host, portToPublish)
	go etcd.PublishService(this.Etcd, this.Host.String(), this.EtcdPath, portToPublish, uint64(this.TTL.Seconds()), this.Timeout)
}

// RunProcessAsDaemon start a child process using a goroutine
func (this *Boot) RunProcessAsDaemon(command string, args []string) {
	go RunProcessAsDaemon(signalChan, command, args)
}

func (this *Boot) RunScript(script string, params map[string]string, loader func(string) ([]byte, error)) {
	RunScript(signalChan, script, params, loader)
}

// WaitForLocalConnection wait until the port/ports exposed are opened
// If no port is specified we use the defined in the constructor
func (this *Boot) WaitForLocalConnection(ports ...string) {
	if len(ports) == 0 {
		Log.Debugf("waiting for a service in the port %v", this.Port)
		WaitForPort("tcp", "127.0.0.1", this.Port, this.Timeout)
	} else {
		// we need to wait for a port different than the default or more than one
		Log.Debugf("waiting for the services in the port/s %v", ports)
		for _, port := range ports {
			WaitForPort("tcp", "127.0.0.1", port, this.Timeout)
		}
	}
}

func (this *Boot) Wait() {
	// wait for exit
	<-signalChan
}
