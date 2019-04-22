package clefui

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

func StartClef(ctx context.Context, binPath string) (*Clef, error) {
	cmd := exec.Command(
		binPath,
		"--keystore", "./tmp/clef/keystore",
		"--configdir", "./tmp/clef/config",
		"--stdio-ui",
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Clef{cmd, stdin, stdout}, nil
}

type Clef struct {
	cmd    *exec.Cmd
	stdin  io.Writer
	stdout io.Reader
}

func (c *Clef) Stop() error {
	log.Info("Stopping Clef")
	if err := c.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.Warn("Error sending signal to Clef", "err", err)
		return c.cmd.Process.Kill()
	}
	waitErr := make(chan error)
	go func() {
		waitErr <- c.cmd.Wait()
	}()
	select {
	case err := <-waitErr:
		return err
	case <-time.After(10 * time.Second):
		log.Warn("Timed out waiting for Clef to stop, killing it")
		return c.cmd.Process.Kill()
	}
}
