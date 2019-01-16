package socks5request

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Request struct {
	conn io.ReadWriter
	bufr *bufio.Reader
}

func New(conn io.ReadWriter) (r *Request) {
	r = &Request{}
	r.conn = conn
	r.bufr = bufio.NewReader(conn)
	return
}

func (self *Request) Do(ip, port int) (err error) {
	// c->s
	_, err = self.conn.Write([]byte{5, 1, 0})
	if err != nil {
		return
	}

	// c<-s
	bs := make([]byte, 2)
	_, err = io.ReadFull(self.conn, bs)
	if err != nil {
		return
	}
	if !(bs[0] == 5 && bs[1] == 0) {
		err = errors.New("socks5 request failed")
		return
	}

	// c->s
	buffer := &bytes.Buffer{}
	binary.Write(buffer, binary.LittleEndian, byte(5))
	binary.Write(buffer, binary.LittleEndian, byte(1))
	binary.Write(buffer, binary.LittleEndian, byte(0))
	binary.Write(buffer, binary.LittleEndian, byte(1))
	binary.Write(buffer, binary.BigEndian, uint32(ip))
	binary.Write(buffer, binary.BigEndian, uint16(port))
	_, err = self.conn.Write(buffer.Bytes())
	if err != nil {
		return
	}

	// c<-s
	bs = make([]byte, 10)
	_, err = io.ReadFull(self.conn, bs)
	if err != nil {
		return
	}
	if !(bs[0] == 5 && bs[1] == 0 && bs[2] == 0 && bs[3] == 1) {
		err = errors.New("socks5 request failed")
		return
	}

	return
}
