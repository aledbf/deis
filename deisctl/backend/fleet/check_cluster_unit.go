package fleet

import (
	"errors"
	"fmt"
	"time"

	"crypto/md5"
	"encoding/hex"
	"github.com/coreos/fleet/ssh"
)

// Status prints the systemd status of target unit(s)
func (c *FleetClient) MD5(name string) (err error) {
	units, err := c.Units(name)
	if err != nil {
		return
	}

	var sshClient *ssh.SSHForwardingClient

	timeout := time.Duration(Flags.SSHTimeout*100) * time.Millisecond

	ms, err := machineState(name)
	if err != nil {
		return err
	}

	if ms == nil {
		machID, err := findUnit(units[0])

		if err != nil {
			return err
		}

		ms, err = machineState(machID)

		if err != nil || ms == nil {
			return err
		}
	}

	addr := ms.PublicIP

	if tun := getTunnelFlag(); tun != "" {
		sshClient, err = ssh.NewTunnelledSSHClient("core", tun, addr, getChecker(), false, timeout)
	} else {
		sshClient, err = ssh.NewSSHClient("core", addr, getChecker(), false, timeout)
	}

	if err != nil {
		return err
	}

	defer sshClient.Close()

	u, err := cAPI.Unit(name)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to get status for unit 1 %s", name))
	}

	if u == nil {
		return errors.New(fmt.Sprintf("Unit %s does not exist 2", name))
	}

	sess, err := sshClient.NewSession()

	cmd := fmt.Sprintf("fleetctl cat %s", name)
	svcContent, err := sess.Output(cmd)

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to get status for unit 3 %s", name))
	}

	md5sum := md5.Sum(svcContent)
	output := hex.EncodeToString(md5sum[:])
	fmt.Printf("%s MD5: %s\n", name, output)
	fmt.Println()
	return
}
