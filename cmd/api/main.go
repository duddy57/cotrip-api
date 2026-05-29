package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/duddy57/cotrip-api/internal/di"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := di.Run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "something went wrong: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "all systems offline, exiting...")
}
