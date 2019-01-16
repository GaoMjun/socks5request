package socks5request

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/GaoMjun/tcpip"

	socks5 "github.com/armon/go-socks5"
)

func TestRequest(t *testing.T) {
	var (
		err      error
		conn     net.Conn
		r        *Request
		s        = "GET / HTTP/1.1\r\nHost: baidu.com\r\n\r\n"
		response *http.Response
		bs       []byte
	)
	defer func() {
		if conn != nil {
			conn.Close()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	go runServer()

	time.Sleep(time.Second * 3)

	conn, err = net.Dial("tcp", "127.0.0.1:1080")
	if err != nil {
		return
	}

	r = New(conn)
	err = r.Do(tcpip.InetAtoN("123.125.115.110"), 80)
	if err != nil {
		return
	}

	fmt.Fprint(conn, s)
	response, err = http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return
	}

	bs, err = httputil.DumpResponse(response, false)
	if err != nil {
		return
	}

	fmt.Println(string(bs))
}

func runServer() {
	var (
		err    error
		config = &socks5.Config{}
		server *socks5.Server
	)
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	server, err = socks5.New(config)
	if err != nil {
		return
	}

	err = server.ListenAndServe("tcp", "127.0.0.1:1080")
}
