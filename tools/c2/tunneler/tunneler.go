package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var mode = [5]string{"ssh", "cmd", "ftp", "sftp", "sh"}

type Command struct {
	program int // 0 for ssh
	args    []string
	cmd     *exec.Cmd
}

func (c *Command) exec() error {
	c.cmd = exec.Command(mode[c.program], c.args...)
	return c.cmd.Run()
}

func main() {
	params := os.Args[0]
	if len(os.Args[:]) < 4 {
		fmt.Println(fmt.Sprintf("usage: %s [mode] [OPTIONS] [LOCAL_IP:]LOCAL_PORT:DESTINATION:DESTINATION_PORT [USER@]SSH_SERVER", params))
		return
	}

	progtype, shellArgs := os.Args[1], os.Args[2:]
	_mode, _ := strconv.Atoi(progtype)
	c := Command{
		program: _mode,
		args:    shellArgs,
	}
	err := c.exec()
	if err != nil {
		fmt.Println(err)
	}
}
