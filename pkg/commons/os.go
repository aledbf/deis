package commons

import (
	"bytes"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/deis/deis/pkg/logger"
	"github.com/progrium/go-basher"
)

const (
	networkWaitTime time.Duration = 5 * time.Second
)

// Getopt return the value of and environment variable or a default
func Getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		logger.Log.Debugf("returning default value \"%s\" for key \"%s\"", dfault, name)
		value = dfault
	}
	return value
}

func StartServiceCommand(signalChan chan os.Signal, command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		logger.Log.Printf("an error ocurred executing command: %v", err)
		signalChan <- syscall.SIGTERM
	}

	err = cmd.Wait()
	logger.Log.Printf("command finished with error: %v", err)
	signalChan <- syscall.SIGTERM
}

func RunBashScript(signalChan chan os.Signal, script string, params map[string]string, loader func(string) ([]byte, error)) {
	bash, _ := basher.NewContext("/bin/bash", false)
	bash.Source(script, loader)
	if params != nil {
		for key, value := range params {
			bash.Export(key, value)
		}
	}

	_, err := bash.Run("main", nil)
	if err != nil {
		logger.Log.Fatal(err)
		signalChan <- syscall.SIGTERM
	}
}

func RunCommand(signalChan chan os.Signal, command string, args []string) string {
	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		logger.Log.Printf("an error ocurred executing command: %v", err)
		signalChan <- syscall.SIGTERM
	}

	return stdout.String()
}

func WaitForLocalConnection(protocol string, testPort string) {
	for {
		_, err := net.DialTimeout(protocol, "127.0.0.1:"+testPort, networkWaitTime)
		if err == nil {
			break
		}
	}
}

// BuildCommandFromString parses a string containing a command and multiple
// arguments and returns a valid tuple to pass to exec.Command
func BuildCommandFromString(input string) (string, []string) {
	command := strings.Split(input, " ")

	if len(command) > 1 {
		return command[0], command[1:]
	}

	return command[0], []string{}
}
