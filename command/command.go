package command

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"time"
)

type Command interface {
	String() string
	ValidateResponse([]byte) ([]byte, error)
}

type Conn interface {
	io.ReadWriter
	SetDeadline(time.Time) error
}

const (
	maxAttempts         = 3
	timeout             = time.Second * 3
	readBufferSizeBytes = 4_096
)

// Send writes the command to the provided Writer, and reads and validates the
// response.
func Send(cmd Command, conn Conn) ([]byte, error) {
	var (
		resp     []byte
		err      error
		attempts int
	)

	for {
		if resp, err = tryRequest(cmd, conn); err != nil {
			attempts++
			log.Printf("error executing command (attempt %d): %s", attempts, err)
			if attempts <= 3 {
				continue
			}
			return nil, fmt.Errorf("error executing command: %s", err)
		}

		return resp, nil
	}
}

func tryRequest(cmd Command, conn Conn) ([]byte, error) {
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("error setting deadline: %s", err)
	}

	p := make([]byte, readBufferSizeBytes)
	_, err := fmt.Fprint(conn, cmd.String())
	if err != nil {
		return nil, fmt.Errorf("error writing to socket: %s", err)
	}

	r := bufio.NewReader(conn)
	n, err := r.Read(p)
	if err != nil {
		return nil, fmt.Errorf("error reading from socket: %s", err)
	}

	return cmd.ValidateResponse(p[:n])
}
