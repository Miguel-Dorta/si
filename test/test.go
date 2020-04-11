package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var (
	tmpDir = filepath.Join(os.TempDir(), "si-test")
	instance1Path = filepath.Join(tmpDir, "instance1")
	instance2Path = filepath.Join(tmpDir, "instance2")
	runningInstances = make([]*exec.Cmd, 1)
)

func main() {
	run("rm", "-Rf", tmpDir)
	run("mkdir", tmpDir)

	run("go", "build", "-o", instance1Path, "instance1.go", "common.go")
	run("go", "build", "-o", instance2Path, "instance2.go", "common.go")

	now := strconv.Itoa(int(time.Now().Unix()))

	instance2, err := createInstance2(instance2Path, now)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := instance2.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error initializing instance2: %s\n", err)
		os.Exit(1)
	}
	runningInstances = append(runningInstances, instance2)

	run(instance1Path, now)
}

func run(s ...string) {
	stdout, stderr := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd := exec.Command(s[0], s[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "error found executing command %v: %s\n>> stdout: %s\n>> stderr: %s\n", s, err, stdout.String(), stderr.String())
	for _, i := range runningInstances {
		i.Process.Kill()
	}
	os.Exit(1)
}

func createInstance2(s ...string) (*exec.Cmd, error) {
	cmd := exec.Command(s[0], s[1:]...)

	stdin, _, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdin pipe: %w", err)
	}
	cmd.Stdin = stdin

	_, stdout, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdout pipe: %w", err)
	}
	cmd.Stdout = stdout

	_, stderr, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stderr pipe: %w", err)
	}
	cmd.Stderr = stderr

	return cmd, nil
}
