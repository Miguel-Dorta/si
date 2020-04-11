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

// Process is an extension of os.Process
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

// StdinPipe creates a pipe to the stdin of the process
func (p *Process) StdinPipe() (io.WriteCloser, error) {
	return os.OpenFile(fmt.Sprintf("/proc/%d/fd/%d", p.Pid, fdStdin), os.O_WRONLY, 0)
}

// StdoutPipe creates a pipe to the stdout of the process
func (p *Process) StdoutPipe() (io.ReadCloser, error) {
	return os.OpenFile(fmt.Sprintf("/proc/%d/fd/%d", p.Pid, fdStdout), os.O_RDONLY, 0)
}

// StderrPipe creates a pipe to the stderr of the process
func (p *Process) StderrPipe() (io.ReadCloser, error) {
	return os.OpenFile(fmt.Sprintf("/proc/%d/fd/%d", p.Pid, fdStderr), os.O_RDONLY, 0)
}
