package latency

import "time"

// Config defines configuration for latency project
type Config struct {
	Address     string        `yaml:"address"`
	DestAddress string        `yaml:"destAddress"`
	Port        int           `yaml:"port"`
	DestPort    int           `yaml:"destPort"`
	Latency     time.Duration `yaml:"latency"`
	QueueSize   int           `yaml:"queueSize"`
}
