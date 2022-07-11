package command

import (
	"bufio"
	"fmt"
	"io"
)

type command interface {
	String() string
	validateResponse([]byte) ([]byte, error)
}

// Send writes the command to the provided Writer, and reads and validates the
// response.
//
// TODO: accept a context.Context and enforce deadline/timeout.
func Send(cmd command, conn io.ReadWriter) ([]byte, error) {
	_, err := fmt.Fprint(conn, cmd.String())
	if err != nil {
		return nil, fmt.Errorf("error writing to socket: %s", err)
	}

	p := make([]byte, 4_096)
	r := bufio.NewReader(conn)
	n, err := r.Read(p)
	if err != nil {
		return nil, fmt.Errorf("error reading from socket: %s", err)
	}

	return cmd.validateResponse(p[:n])
}
