package si

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
)

const (
	fdStdin = iota
	fdStdout
	fdStderr
)

type Process struct {
	*os.Process
}

func findProcess(pid int) (*Process, error) {
	if err := unix.Kill(pid, unix.Signal(0)); err != nil {
		return nil, err
	}
	p, err := os.FindProcess(pid)
	return &Process{p}, err
}

func (p *Process) StdinPipe() (io.WriteCloser, error) {
	return os.OpenFile(fmt.Sprintf("/proc/%d/fd/%d", p.Pid, fdStdin), os.O_WRONLY|os.O_APPEND, 0)
}

func (p *Process) StdoutPipe() (io.ReadCloser, error) {
	return os.OpenFile(fmt.Sprintf("/proc/%d/fd/%d", p.Pid, fdStdout), os.O_RDONLY|os.O_APPEND, 0)
}

func (p *Process) StderrPipe() (io.ReadCloser, error) {
	return os.OpenFile(fmt.Sprintf("/proc/%d/fd/%d", p.Pid, fdStderr), os.O_RDONLY|os.O_APPEND, 0)
}
