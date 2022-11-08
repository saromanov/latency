package latency

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type Latency struct {
	cfg                         Config
	tcpAddress, resolvedAddress *net.TCPAddr
	listener                    *net.TCPListener
	cancelFunc                  context.CancelFunc
}

// New provides definition of the Latency
func New(cfg Config) *Latency {
	return &Latency{
		cfg: cfg,
	}
}

// Init provides initialization of Latency
func (l *Latency) Init(ctx context.Context) error {
	log := logrus.WithContext(ctx)
	tcpAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", l.cfg.Port))
	if err != nil {
		return fmt.Errorf("Error resolving local address: %s", err)
	}
	resolvedTCPAddress, err := net.ResolveTCPAddr("tcp", l.cfg.Address)
	if err != nil {
		return fmt.Errorf("Error resolving destination address: %s", err)
	}
	log.Info("Latency was initialized")
	l.tcpAddress = tcpAddress
	l.resolvedAddress = resolvedTCPAddress
	return nil
}

// Start provides starting of the Latency server and connect to
func (s *Latency) Start() error {
	listener, err := net.ListenTCP("tcp", s.tcpAddress)
	if err != nil {
		return fmt.Errorf("Error starting TCP listener: %s", err)
	}
	if listener == nil {
		return fmt.Errorf("listener is not defined try it again")
	}
	s.listener = listener

	ctx, cancel := context.WithCancel(context.Background())
	if err := s.start(ctx, cancel); err != nil {
		return err
	}
	s.cancelFunc = cancel
	return nil
}

// start provides starting of accepting loop of connections
func (s *Latency) start(ctx context.Context, cancel context.CancelFunc) error {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			return err
		}

		c := NewConnection(conn)
		if err := c.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}
// stop poriding sriopping of Latency
func (s *Latency) Stop(ctx context.Context) error {
	log := logrus.WithContext(ctx)
	log.Info("Stopping of the Latency")
	s.cancelFunc()
	return nil
}
