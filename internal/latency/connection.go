package latency

import (
	"context"
	"net"
)

type Connection struct {
	conn *net.TCPConn
	bufferSize int
}

func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Start(ctx context.Context) error {
	buffer := make([]byte, c.bufferSize)
		for {
			
		}
	return nil
}