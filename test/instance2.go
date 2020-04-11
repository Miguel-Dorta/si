package main

import (
	"fmt"
	"github.com/Miguel-Dorta/si"
	"os"
	"strconv"
	"time"
)

func main() {
	instanceName = "instance2"

	if err := si.Register(alias); err != nil {
		failf("error registering process: %s", err)
	}

	// Wait to instance1 to start and register
	time.Sleep(time.Second*2)

	s, err := readString(os.Stdin)
	if err != nil {
		failf("error reading from stdin: %s", err)
	}
	if s != strconv.Itoa(os.Getpid()) {
		failf("data from stdin: expected \"%d\", got \"%s\"", os.Getpid(), s)
	}

	fmt.Println("stdout")
	_, _ = fmt.Fprintln(os.Stderr, "stderr")

	time.Sleep(time.Minute)
	fail("timeout, exiting process")
}
