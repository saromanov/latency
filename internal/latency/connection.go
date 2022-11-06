package latency

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
)

type Connection struct {
	conn, writeConn *net.TCPConn
	bufferSize      int
	done            chan error
}

func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		conn: conn,
		done: make(chan error),
	}
}

func (c *Connection) Start(ctx context.Context) error {
	log := logrus.WithContext(ctx)
	go c.readFromBuffer(ctx)
	for {
		select {
		case err := <-c.done:
			log.WithError(err).Error("error on done method")
			return err
		}
	}
}

func (c *Connection) readFromBuffer(ctx context.Context) error {
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
