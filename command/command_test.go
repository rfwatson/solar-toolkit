package command_test

import (
	"errors"
	"testing"
	"time"

	"git.netflux.io/rob/solar-toolkit/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type readResult struct {
	p   []byte
	err error
}

type mockConn struct {
	readResults []readResult
}

func (c *mockConn) Read(p []byte) (int, error) {
	var result readResult
	result, c.readResults = c.readResults[0], c.readResults[1:]
	return copy(p, result.p), result.err
}

func (c *mockConn) Write(p []byte) (int, error) { return 0, nil }
func (c *mockConn) SetDeadline(time.Time) error { return nil }

type mockCommand struct{}

func (cmd *mockCommand) String() string                            { return "baz" }
func (cmd *mockCommand) ValidateResponse(p []byte) ([]byte, error) { return p, nil }

func TestSendWithOneRetry(t *testing.T) {
	var cmd mockCommand
	conn := mockConn{
		readResults: []readResult{
			{err: errors.New("i/o timeout")},
			{p: []byte("bar"), err: nil},
		},
	}

	resp, err := command.Send(&cmd, &conn)
	require.NoError(t, err)
	assert.Equal(t, []byte("bar"), resp)
}

func TestSendFail(t *testing.T) {
	var cmd mockCommand
	conn := mockConn{
		readResults: []readResult{
			{err: errors.New("i/o timeout 1")},
			{err: errors.New("i/o timeout 2")},
			{err: errors.New("i/o timeout 3")},
			{err: errors.New("i/o timeout 4")},
		},
	}

	_, err := command.Send(&cmd, &conn)
	assert.EqualError(t, err, "error executing command: error reading from socket: i/o timeout 4")
}
