package commons

import (
	"fmt"
	"os"
	"time"

	"github.com/deis/deis/pkg/logger"
)

// WaitForInitialConfd wait until the compilation of the templates is correct
func WaitForInitialConfd(signalChan chan os.Signal, etcd string, timeout time.Duration) {
	for {
		cmdAsString := fmt.Sprintf("confd -onetime -node %s -confdir /app", etcd)
		cmd, args := BuildCommandFromString(cmdAsString)
		logger.Log.Info("waiting for confd to write initial templates...")
		err := RunCommand(signalChan, cmd, args, false)
		if err == nil {
			break
		}

		time.Sleep(timeout)
	}
}

// LaunchConfd launch confd as a daemon process.
func LaunchConfd(signalChan chan os.Signal, etcd string) {
	cmdAsString := fmt.Sprintf("confd -node %s -confdir /app --interval 5 --quiet --watch", etcd)
	cmd, args := BuildCommandFromString(cmdAsString)
	go RunProcessAsDaemon(signalChan, cmd, args)
}
