package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/log"
	"github.com/lmars/clefui"
	isatty "github.com/mattn/go-isatty"
)

func main() {
	// configure logging
	var format log.Format
	if isatty.IsTerminal(os.Stdout.Fd()) {
		format = log.TerminalFormat(true)
	} else {
		format = log.LogfmtFormat()
	}
	log.Root().SetHandler(log.StreamHandler(os.Stdout, format))

	// shutdown gracefully on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Info("Received signal, exiting...")
	}()

	// run the Clef UI
	if err := clefui.Run(ctx); err != nil {
		log.Crit("Error running Clef UI", "err", err)
	}
}
