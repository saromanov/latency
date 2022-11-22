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
	destAddress     *net.TCPAddr
	bufferSize      int
	done            chan error
	delayedRequest  chan delayedRequest
}

type delayedRequest struct {
	data  []byte
	delay time.Duration
}

func NewConnection(conn *net.TCPConn, destAddress *net.TCPAddr) *Connection {
	return &Connection{
		conn:           conn,
		destAddress:    destAddress,
		done:           make(chan error),
		delayedRequest: make(chan delayedRequest),
	}
}

func (c *Connection) Start(ctx context.Context) error {
	log := logrus.WithContext(ctx)
	destConn, err := net.DialTCP("tcp", nil, c.destAddress)
	if err != nil {
		return fmt.Errorf("Error dialing remote address: %s", err)
	}
	c.writeConn = destConn
	go c.readFromSrc(ctx)
	go c.handleDelayedRequests(ctx)
	for {
		select {
		case err := <-c.done:
			log.WithError(err).Error("error on done method")
			return err
		case <- ctx.Done():
			return fmt.Errorf("context was cancelled")
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
		data := <-c.delayedRequest
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
