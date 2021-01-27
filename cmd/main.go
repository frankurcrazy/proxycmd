package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/frankurcrazy/proxycmd"
)

func readHTTPHeader(c proxycmd.BufferedConn) (string, error) {
	const (
		Idle         int = iota
		HalfDelim        = iota
		FirstDelim       = iota
		OneHalfDelim     = iota
	)

	buf := &bytes.Buffer{}
	state := Idle

	for {
		b, err := c.ReadByte()
		if err != nil {
			return "", err
		}

		buf.WriteByte(b)

		switch b {
		case 13:
			if state == Idle {
				state = HalfDelim
			} else if state == FirstDelim {
				state = OneHalfDelim
			} else {
				state = Idle
			}
		case 10:
			if state == OneHalfDelim {
				return buf.String(), nil
			} else if state == HalfDelim {
				state = FirstDelim
			} else {
				state = Idle
			}
		default:
			state = Idle
		}
	}
}

func connectProxy(u string) (net.Conn, error) {
	parsedURL, err := url.Parse(u)
	conn, err := net.Dial("tcp", parsedURL.Host)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to proxy server: %s", err)

	}

	return conn, nil
}

func connectSSH(conn net.Conn, s string) (net.Conn, error) {
	bConn := proxycmd.NewBufferedConn(conn)
	_, err := bConn.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", s)))
	if err != nil {
		return nil, err
	}

	_, err = bConn.Peek(1)
	if err != nil {
		return nil, err
	}

	header, err := readHTTPHeader(bConn)
	if err != nil {
		return nil, err
	}

	if strings.Contains(header, "HTTP/1.1 502") {
		return nil, fmt.Errorf("Error 502")
	}

	return bConn, nil
}

func serviceSSHConn(conn net.Conn, wg *sync.WaitGroup, doneCh chan int) {
	defer wg.Done()
	defer conn.Close()
	buf := make([]byte, 2048)

	// stdin reader
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		buf := make([]byte, 2048)

		for {
			n, err := reader.Read(buf)
			if err != nil {
				return
			}

			conn.Write(buf[:n])
		}
	}()

	// conn reader
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		os.Stdout.Write(buf[:n])
	}
}

func main() {
	proxyStr := flag.String("proxy", "http://10.107.60.10:31280", "HTTP Proxy to use")
	sshServer := flag.String("ssh-server", "", "SSH server and port to connect")
	flag.Parse()
	logger := log.New(os.Stderr, "[proxycmd] ", log.Ldate|log.Lshortfile)

	if len(*sshServer) == 0 {
		flag.Usage()
		return
	}
	if len(*proxyStr) == 0 {
		flag.Usage()
		return
	}

	logger.Printf("connecting to %s via %s\n", *sshServer, *proxyStr)

	wg := &sync.WaitGroup{}
	doneCh := make(chan int)

	proxyConn, err := connectProxy(*proxyStr)
	if err != nil {
		logger.Fatalf("failed connecting to proxy: %s\n", err)
	}

	sshConn, err := connectSSH(proxyConn, *sshServer)
	if err != nil {
		logger.Fatalf("failed connecting to ssh server: %s\n", err)
	}

	wg.Add(1)
	go serviceSSHConn(sshConn, wg, doneCh)
	wg.Wait()
}
