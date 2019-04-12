package signer

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/log"
)

var ErrShutdown = errors.New("signer: Signer stopped")

func New() *Signer {
	return &Signer{}
}

type Signer struct {
}

func (s *Signer) Run(ctx context.Context) error {
	log.Info("Starting the signer UI")
	<-ctx.Done()
	log.Info("Stopping the signer UI")
	return ErrShutdown
}
