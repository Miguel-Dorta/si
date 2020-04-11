package main

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/si"
	"os"
	"time"
)

func main() {
	instanceName = "instance1"

	// Wait for instance2 to start and register
	time.Sleep(time.Second)

	// TEST Register
	if err := si.Register(alias); err == nil {
		fail("no error when registering same alias")
	}
	pass("Register")

	// TEST Find non-existent
	p, err := si.Find("non-existent")
	if err != nil {
		failf("error finding non-existent: %s", err)
	}
	if p != nil {
		failf("process found looking for non-existent alias (PID=%d)", p.Pid)
	}
	pass("Find non-existing")

	// TEST Find existent
	p, err = si.Find(alias)
	if err != nil {
		failf("error finding existing alias: %s", err)
	}
	if p == nil {
		fail("could not find running instance")
	}
	if err := testProcess(p); err != nil {
		fail(err.Error())
	}
	pass("Find existing")

	if err := p.Kill(); err != nil {
		failf("error killing process: %s", err)
	}
}

func testProcess(p *si.Process) error {
	stdin, err := p.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdin pipe: %s", err)
	}
	stdout, err := p.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %s", err)
	}
	stderr, err := p.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %s", err)
	}

	if _, err := fmt.Fprintln(stdin, p.Pid); err != nil {
		return fmt.Errorf("error writing to stdin pipe: %s", err)
	}
	if s, err := readString(stdout); err != nil {
		return fmt.Errorf("error reading from stdout: %s", err)
	} else if s != "stdout" {
		return fmt.Errorf("expecting \"stdout\", found \"%s\"", s)
	}
	if s, err := readString(stderr); err != nil {
		return fmt.Errorf("error reading from stderr: %s", err)
	} else if s != "stderr" {
		fmt.Fprintln(os.Stderr, s)
		return errors.New("expecting \"stderr\", found other thing")
	}
	return nil
}

func pass(testName string) {
	fmt.Printf("PASS: %s test\n", testName)
}
