package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"collector/processor"
	"collector/publisher"
	"collector/sensor"
)

func main() {
	opts := parseFlags(os.Args[1:])

	slog.Info("Starting sensor node", "name", opts.Name, "destination", opts.DstAddress)

	publ := new(publisher.Publisher)
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
	DstAddress string
	Rate       uint
}

func parseFlags(args []string) (opts Options) {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&opts.Name, "name", "node", "Name of the sensor")
	fs.StringVar(&opts.DstAddress, "dst", "", "Destination address where to send data (required)")

	fs.UintVar(&opts.Rate, "rate", 100, "Number of messages per second to send")

	if err := fs.Parse(args); err != nil {
		panic(err) // should never happend
	}

	if opts.DstAddress == "" {
		fmt.Println("Required option `dst` is not specified")
		fs.Usage()
		os.Exit(2)
	}

	return opts
}
