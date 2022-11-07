package latency

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type Connection struct {
	conn, writeConn *net.TCPConn
	bufferSize      int
	done            chan error
	delayedRequest  chan delayedRequest
}

type delayedRequest struct {
	data  []byte
	delay time.Duration
}

func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		conn:           conn,
		done:           make(chan error),
		delayedRequest: make(chan delayedRequest),
	}
}

func (c *Connection) Start(ctx context.Context) error {
	log := logrus.WithContext(ctx)
	go c.readFromSrc(ctx)
	for {
		select {
		case err := <-c.done:
			log.WithError(err).Error("error on done method")
			return err
		}
	}
}

// handleDest provides handling of destination buffer
func (c *Connection) handleDest(ctx context.Context) error {
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

// readFromSrc provides reading from src host
// also, it sends to delayed queue
func (c *Connection) readFromSrc(ctx context.Context) error {
	for {
		buffer := make([]byte, c.bufferSize)
		data, err := c.writeConn.Read(buffer)
		if err != nil {
			return err
		}

		d := delayedRequest{
			data: buffer[:data],
		}
		c.delayedRequest <- d

	}
	return nil
}

// handling of delayed requests
func (c *Connection) handleDelayedRequests(ctx context.Context) error {
	log := logrus.WithContext(ctx)
	for {
		data := <- c.delayedRequest
		if data.delay.Seconds() == 0 {
			log.Info("seconds is zero. There is not delay")
		} else {
			time.Sleep(data.delay)
		}
		_, err := c.writeConn.Write(data.data)
		if err != nil {
			c.done <- fmt.Errorf("Error writing to queue: %s", err)
			return err
		}
	}
	return nil
}
