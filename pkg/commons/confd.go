package commons

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/deis/deis/pkg/logger"
)

// WaitForInitialConfd wait until the compilation of the templates is correct
func WaitForInitialConfd(etcd string, timeout time.Duration) {
	for {
		var buffer bytes.Buffer
		output := bufio.NewWriter(&buffer)
		cmd := exec.Command("confd", "-onetime", "-node", etcd, "-confdir", "/app")
		cmd.Stdout = output
		cmd.Stderr = output
		err := cmd.Run()
		output.Flush()
		if err == nil {
			break
		}

		logger.Log.Info("waiting for confd to write initial templates...")
		logger.Log.Debugf("%v", buffer.String())
		time.Sleep(timeout)
	}
}

// LaunchConfd Launch confd as a daemon process
func LaunchConfd(signalChan chan os.Signal, etcd string) {
	cmd := exec.Command("confd", "-node", etcd, "-confdir", "/app", "--interval", "5", "--quiet", "--watch")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.Log.Errorf("confd terminated by error: %v", err)
		signalChan <- syscall.SIGTERM
	}

	err := cmd.Wait()
	logger.Log.Printf("confd terminated with error: %v", err)
	signalChan <- syscall.SIGTERM
}
