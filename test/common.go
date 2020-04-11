package main

import (
	"fmt"
	"github.com/Miguel-Dorta/si"
	"io"
	"os"
)

var (
	alias string
	instanceName string
)

func init() {
	si.Dir = os.TempDir()
	alias = "si_test_" + os.Args[1]
}

func readString(r io.Reader) (string, error) {
	b := make([]byte, 10)
	n, err := r.Read(b)
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return "", err
	}

	if n == 0 {
		n++
	}
	return string(b[:n-1]), nil
}

func failf(format string, a ...interface{}) {
	fail(fmt.Sprintf(format, a...))
}

func fail(s string) {
	fmt.Printf("FAIL [%s]: %s\n", instanceName, s)
	os.Exit(1)
}
