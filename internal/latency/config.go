package latency

import "time"

// Config defines configuration for latency project
type Config struct {
	Address     string    `yaml:"address"`
	DestAddress string    `yaml:"destAddress"`
	Port        int       `yaml:"port"`
	Latency     time.Time `yaml:"latency"`
	QueueSize   int       `yaml:"queueSize"`
}
