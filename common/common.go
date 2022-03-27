package common

import (
	"fmt"
	"strconv"
)

const (
	SocketPath = "/tmp/ipc-unix.sock"
	Protocol   = "unix"     // "tcp"
	SocketAddr = SocketPath //"localhost:8080"
)

type Data struct {
	Stdout []byte
	Stderr []byte
	Msg    string
}

func (d *Data) Reset() {
	d.Stdout = nil
	d.Stderr = nil
	d.Msg = ""
}

func (d Data) String() string {
	return fmt.Sprintf(
		`stdout: %s
stderr: %s
msg   : %s
	`, strconv.Quote(string(d.Stdout)), strconv.Quote(string(d.Stderr)), d.Msg)
}
