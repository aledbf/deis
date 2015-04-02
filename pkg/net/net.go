package net

import (
	"net"
	"time"
)

// WaitForPort wait for successful network connection
func WaitForPort(proto string, ip string, port string, timeout time.Duration) error {
	for {
		con, err := net.DialTimeout(proto, ip+":"+port, timeout)
		if err == nil {
			con.Close()
			break
		}
	}

	return nil
}
