package boot

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-etcd/etcd"

	"github.com/deis/deis/pkg/commons"
	"github.com/deis/deis/pkg/logger"
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
	logger.Log.Info("starting deis component...")

	host := commons.Getopt("HOST", "127.0.0.1")
	etcdPort := commons.Getopt("ETCD_PORT", "4001")

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
	commons.WaitForInitialConfd(signalChan, this.Host.String()+":"+this.EtcdPort, this.Timeout)

	// spawn confd in the background to update services based on etcd changes
	go commons.LaunchConfd(signalChan, this.Host.String()+":"+this.EtcdPort)
}

// Publish publish information about the relevant process running in the boot
// process in etcd using specified path and port/s
func (this *Boot) Publish(port ...string) {
	portToPublish := this.Port
	// If we specify a custom port we use that one
	if len(port) != 0 {
		portToPublish = port[1]
	}
	logger.Log.Info("starting periodic publication in etcd...")
	logger.Log.Debugf("etcd publication path %s, host %s and port %s", this.EtcdPath, this.Host, portToPublish)
	go commons.PublishService(this.Etcd, this.Host.String(), this.EtcdPath, portToPublish, uint64(this.TTL.Seconds()), this.Timeout)
}

// RunProcessAsDaemon start a child process using a goroutine
func (this *Boot) RunProcessAsDaemon(command string, args []string) {
	go commons.RunProcessAsDaemon(signalChan, command, args)
}

func (this *Boot) RunScript(script string, params map[string]string, loader func(string) ([]byte, error)) {
	commons.RunScript(signalChan, script, params, loader)
}

// WaitForLocalConnection wait until the port/ports exposed are opened
// If no port is specified we use the defined in the constructor
func (this *Boot) WaitForLocalConnection(ports ...string) {
	if len(ports) == 0 {
		logger.Log.Debugf("waiting for a service in the port %v", this.Port)
		commons.WaitForPort("tcp", "127.0.0.1", this.Port, this.Timeout)
	} else {
		// we need to wait for a port different than the default or more than one
		logger.Log.Debugf("waiting for the services in the port/s [%v]", ports)
		for _, port := range ports {
			commons.WaitForPort("tcp", "127.0.0.1", port, this.Timeout)
		}
	}
}

func (this *Boot) Wait() {
	// wait for exit
	<-signalChan
}
