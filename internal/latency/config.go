package latency

import "time"

// Config defines configuration for latency project
type Config struct {
	Address     string
	DestAddress string
	Port        int
	DestPort    int
	Latency     time.Duration
	QueueSize   int
	BufferSize  int
}
