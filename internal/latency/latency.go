package latency

import (
	"fmt"
	"net"
)

type Latency struct {
	cfg                         Config
	tcpAddress, resolvedAddress *net.TCPAddr
	listener                    *net.TCPListener
}

// New provides definition of the Latency
func New(cfg Config) *Latency {
	return &Latency{
		cfg: cfg,
	}
}

// Init provides initialization of Latency
func (l *Latency) Init() error {
	tcpAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", l.cfg.Port))
	if err != nil {
		return fmt.Errorf("Error resolving local address: %s", err)
	}
	resolvedTCPAddress, err := net.ResolveTCPAddr("tcp", l.cfg.Address)
	if err != nil {
		return fmt.Errorf("Error resolving destination address: %s", err)
	}
	l.tcpAddress = tcpAddress
	l.resolvedAddress = resolvedTCPAddress
	return nil
}

// Start provides starting of the Latency
func (s *Latency) Start() error {
	listener, err := net.ListenTCP("tcp", s.tcpAddress)
	if err != nil {
		return fmt.Errorf("Error starting TCP listener: %s", err)
	}
	s.listener = listener

	return nil
}
