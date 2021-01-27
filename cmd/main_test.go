package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testBufConn struct {
	buf *bytes.Buffer
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (tbc *testBufConn) Read(b []byte) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (tbc *testBufConn) Write(b []byte) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (tbc *testBufConn) Close() error {
	return nil
}

// LocalAddr returns the local network address.
func (tbc *testBufConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

// RemoteAddr returns the remote network address.
func (tbc *testBufConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (tbc *testBufConn) Discard(n int) (int, error) {
	panic("not implemented")
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
//
// Note that if a TCP connection has keep-alive turned on,
// which is the default unless overridden by Dialer.KeepAlive
// or ListenConfig.KeepAlive, then a keep-alive failure may
// also return a timeout error. On Unix systems a keep-alive
// failure on I/O can be detected using
// errors.Is(err, syscall.ETIMEDOUT).
func (tbc *testBufConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (tbc *testBufConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (tbc *testBufConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (tbc *testBufConn) Peek(n int) ([]byte, error) {
	return tbc.buf.Bytes()[:n], nil
}

func (tbc *testBufConn) Buffered() (n int) {
	return tbc.buf.Len()
}

func (tbc *testBufConn) ReadByte() (byte, error) {
	return tbc.buf.ReadByte()
}

func (tbc *testBufConn) UnreadByte() error {
	return tbc.buf.UnreadByte()
}

func TestReadHTTPHeader(t *testing.T) {
	testbuf := bytes.NewBuffer([]byte("HTTP/1.1 200 Connection established\r\nVia: 1.1 SRVTMG01\r\nConnection: Keep-Alive\r\nProxy-Connection: Keep-Alive\r\n\r\nThis is body"))

	bufConn := &testBufConn{
		buf: testbuf,
	}

	res, err := readHTTPHeader(bufConn)
	assert.Nil(t, err)
	assert.False(t, strings.Contains(res, "This is body"))
	assert.True(t, strings.Contains(res, "HTTP/1.1"))
	assert.True(t, strings.Contains(res, "Keep-Alive\r\n\r\n"))
}
