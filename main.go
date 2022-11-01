package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/saromanov/latency/internal/latency"
)

// parseArgs provides parsing of argumenys
func parseArgs(args []string) (latency.Config, error) {
	var app = kingpin.New("latency", "Proxy for simulating network latency")
	var (
		address = app.Flag("address", "Address of the host").Default("localhost").String()
		destAddress = app.Flag("destAddress", "Address of the destanation host").String()
		port = app.Flag("port", "Port").Default("8080").Int()
	)

	if _, err := app.Parse(args); err != nil {
		return latency.Config{}, err
	}

	return latency.Config{
		Address: *address,
		Port: *port,
		DestAddress: *destAddress,
	}, nil
}
func main() {
	data, err := parseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}

	l := latency.New(data)

	if err := l.Init(); err != nil {
		panic(err)
	}
}
