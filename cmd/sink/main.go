package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"collector/handler"
	"collector/writer"
)

func main() {
	opts := parseFlags(os.Args[1:])

	f, err := os.Create(opts.OutPath)
	if err != nil {
		panic(fmt.Errorf("output file can't be created: %w", err))
	}

	defer f.Close()

	bw := writer.New(f, int(opts.BufferSize), opts.FlushInterval)
	defer bw.Close()

	srv := &http.Server{
		Addr:         opts.ListenAddr,
		Handler:      handler.NewLimiter(handler.WriteTo(bw), int64(opts.RateLimit), time.Second),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer cancel() // cancel main context when Shutdown is finished

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		sig := <-ch

		slog.Info("Stopping sever", "signal", sig)

		shCtx, shCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shCancel()

		if err := srv.Shutdown(shCtx); err != nil { // calling Shutdown with timeout
			slog.Warn("Connections closing error on shutdown", "error", err)
		}
	}()

	slog.Info("Start listening", "address", opts.ListenAddr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	<-ctx.Done() // waiting when all connections are closed
}

type Options struct {
	ListenAddr    string
	OutPath       string
	BufferSize    uint
	FlushInterval time.Duration
	RateLimit     uint
}

func parseFlags(args []string) (opts Options) {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&opts.ListenAddr, "listen", ":8080", "`address` to listen on")
	fs.StringVar(&opts.OutPath, "out", "", "Path to the output file (required)")

	fs.UintVar(&opts.BufferSize, "buffer", 1_000_000, "Size of buffer in bytes")
	fs.UintVar(&opts.RateLimit, "rate", 50_000, "Maximum allowed input flow rate in bytes/sec")

	fs.DurationVar(&opts.FlushInterval, "interval", 100*time.Millisecond, "Buffer flush interval")

	if err := fs.Parse(args); err != nil {
		panic(err) // should never happend
	}

	if opts.OutPath == "" {
		fmt.Println("Required option `out` is not specified")
		fs.Usage()
		os.Exit(2)
	}

	return opts
}
