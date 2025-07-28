package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dvachaiev/collector/processor"
	"github.com/dvachaiev/collector/publisher"
	"github.com/dvachaiev/collector/sensor"
)

func main() {
	opts := parseFlags(os.Args[1:])

	slog.Info("Starting sensor node", "name", opts.Name, "destination", opts.DstURL)

	publ := publisher.New(opts.DstURL, int(opts.BufferSize))
	defer publ.Close()

	proc := processor.New("sensor1", new(sensor.Incremental), int(opts.Rate), publ)
	defer proc.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	sig := <-ch

	slog.Info("Stopping node", "signal", sig)
}

type Options struct {
	Name       string
	DstURL     string
	Rate       uint
	BufferSize uint
}

func parseFlags(args []string) (opts Options) {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	var addr string

	fs.StringVar(&opts.Name, "name", "node", "Name of the sensor")
	fs.StringVar(&addr, "dst", "", "Destination address where to send data (required)")

	fs.UintVar(&opts.Rate, "rate", 100, "Number of messages per second to send")
	fs.UintVar(&opts.BufferSize, "buffer", 100_000, "Buffer size in bytes to cache data")

	if err := fs.Parse(args); err != nil {
		panic(err) // should never happend
	}

	if addr == "" {
		fmt.Println("Required option `dst` is not specified")
		fs.Usage()
		os.Exit(2)
	}

	opts.DstURL = fmt.Sprintf("http://%v/", addr)

	return opts
}
