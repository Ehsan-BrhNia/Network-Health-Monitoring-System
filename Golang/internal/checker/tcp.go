package checker

import (
	"fmt"
	"net"
	"time"
)

func TCPCheck(
	host string,
	port int,
	timeout time.Duration,
) error {

	conn, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%d", host, port),
		timeout,
	)

	if err != nil {
		return err
	}

	defer conn.Close()

	return nil
}
