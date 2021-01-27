package proxycmd

import (
	"bufio"
	"net"
	"time"
)

// BufferedConn -
type BufferedConn interface {
	net.Conn
	Peek(n int) ([]byte, error)
	Buffered() (n int)
	ReadByte() (byte, error)
	UnreadByte() error
	Discard(n int) (int, error)
}

type bConn struct {
	reader *bufio.Reader
	conn   net.Conn
}

func NewBufferedConn(c net.Conn) BufferedConn {
	return &bConn{
		bufio.NewReader(c),
		c,
	}
}

func (b *bConn) Discard(n int) (int, error) {
	return b.reader.Discard(n)
}

func (b *bConn) Peek(n int) ([]byte, error) {
	return b.reader.Peek(n)
}

func (b *bConn) Buffered() (n int) {
	return b.reader.Buffered()
}

func (b *bConn) ReadByte() (byte, error) {
	return b.reader.ReadByte()
}

func (b *bConn) UnreadByte() error {
	return b.reader.UnreadByte()
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (b *bConn) Read(d []byte) (n int, err error) {
	return b.reader.Read(d)
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (b *bConn) Write(d []byte) (n int, err error) {
	return b.conn.Write(d)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (b *bConn) Close() error {
	return b.conn.Close()
}

// LocalAddr returns the local network address.
func (b *bConn) LocalAddr() net.Addr {
	return b.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (b *bConn) RemoteAddr() net.Addr {
	return b.conn.RemoteAddr()
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
func (b *bConn) SetDeadline(t time.Time) error {
	return b.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (b *bConn) SetReadDeadline(t time.Time) error {
	return b.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (b *bConn) SetWriteDeadline(t time.Time) error {
	return b.conn.SetWriteDeadline(t)
}
