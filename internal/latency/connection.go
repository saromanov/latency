package latency

import (
	"context"
	"net"
)

type Connection struct {
	conn, writeConn *net.TCPConn
	bufferSize      int
}

func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Start(ctx context.Context) error {
	buffer := make([]byte, c.bufferSize)
	for {
		data, err := c.conn.Read(buffer)
		if err != nil {
			return err
		}
		data, err = c.conn.Write(buffer[:data])
		if err != nil {
			return err
		}
	}
	return nil
}
