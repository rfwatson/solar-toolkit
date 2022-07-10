package command

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

type command interface {
	String() string
	validateResponse([]byte) error
}

func Send(cmd command, conn io.ReadWriter) ([]byte, error) {
	_, err := fmt.Fprint(conn, cmd.String())
	if err != nil {
		return nil, fmt.Errorf("error writing to socket: %s", err)
	}

	log.Printf("sent data to socket: %X", cmd)

	p := make([]byte, 4_096)
	r := bufio.NewReader(conn)
	n, err := r.Read(p)
	if err != nil {
		return nil, fmt.Errorf("error reading from socket: %s", err)
	}
	p = p[:n]

	if err := cmd.validateResponse(p); err != nil {
		return nil, fmt.Errorf("error validating response: %s", err)
	}

	return p, nil
}
