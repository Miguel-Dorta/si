package si

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const ext = ".pid"

var (
	Dir = "/run"
	ErrOtherInstanceRunning = errors.New("there is another instance of the program running")
)

func Find(alias string) (*Process, error) {
	f, err := os.Open(getPath(alias))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error opening pid file: %w", err)
	}
	defer f.Close()

	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB); err == nil {
		if err := unix.Flock(int(f.Fd()), unix.LOCK_UN); err != nil {
			return nil, fmt.Errorf("error unlocking pid file: %w", err)
		}
		return nil, nil
	} else if err != unix.EWOULDBLOCK {
		return nil, fmt.Errorf("error locking pid file: %w", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading pid file: %w", err)
	}
	data = bytes.TrimSpace(data)

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return nil, fmt.Errorf("error parsing pid (%s): %w", string(data), err)
	}

	p, err := findProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("error finding process with pid %d: %w", pid, err)
	}
	return p, err
}

func Register(alias string) error {
	p, err := Find(alias)
	if err != nil {
		return err
	}
	if p != nil {
		return ErrOtherInstanceRunning
	}

	f, err := os.OpenFile(getPath(alias), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error creating pid file: %w", err)
	}
	defer f.Close()

	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		if err == unix.EWOULDBLOCK {
			return ErrOtherInstanceRunning
		}
		return fmt.Errorf("error locking pid file: %w", err)
	}

	if _, err := f.WriteString(strconv.Itoa(os.Getpid()) + "\n"); err != nil {
		return fmt.Errorf("error writing to pid file: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("error closing pid file: %w", err)
	}
	return nil
}

func getPath(alias string) string {
	return filepath.Join(Dir, alias + ext)
}
