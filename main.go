package main

import (
	"context"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/saromanov/latency/internal/latency"
)

// parseArgs provides parsing of argumenys
func parseArgs(args []string) (latency.Config, error) {
	var app = kingpin.New("latency", "Proxy for simulating network latency")
	var (
		latencyDur = app.Flag("latency", "Duration for latency").Duration()
		address = app.Flag("address", "Address of the host").Default("localhost").String()
		destAddress = app.Flag("destAddress", "Address of the destanation host").String()
		port = app.Flag("port", "Port").Default("8080").Int()
		destPort = app.Flag("destPort", "DestPort").Int()
	)

	if _, err := app.Parse(args); err != nil {
		return latency.Config{}, err
	}

	return latency.Config{
		Address: *address,
		Port: *port,
		DestAddress: *destAddress,
		DestPort: *destPort,
		Latency: *latencyDur,
	}, nil
}
func main() {
	data, err := parseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	l := latency.New(data)

	if err := l.Init(ctx); err != nil {
		panic(err)
	}

	if err := l.Start(ctx); err != nil {
		panic(err)
	}
}
