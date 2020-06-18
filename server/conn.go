package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	IdleTimeout time.Duration
	MaxReadBuffer int64
}

func (connInterface *Conn) Write(p []byte) error {
	connInterface.updateDeadline()
	_, err := connInterface.Conn.Write(p)
	if err != nil {
		return err
	}
	return nil
}

func (connInterface *Conn) Read(b []byte) (int, error) {
	connInterface.updateDeadline()
	reader := io.LimitReader(connInterface.Conn, connInterface.MaxReadBuffer)
	return reader.Read(b)
}

func (connInterface *Conn) updateDeadline() {
	idleDeadline := time.Now().Add(connInterface.IdleTimeout)
	_ = connInterface.Conn.SetDeadline(idleDeadline)
}

func (connInterface *Conn) GetReader() io.Reader {
	return io.LimitReader(connInterface.Conn, connInterface.MaxReadBuffer)
}

func (connInterface *Conn) GetRemoteAddr() string {
	return connInterface.Conn.RemoteAddr().String()
}


type InterfaceData struct {
	id int
	conn *Conn
}

func (i *InterfaceData) GetClientResponse() (int, []byte) {
	var buff []byte
	var response Response
	_, err := i.conn.Read(buff)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(buff, &response)
	return response.Code, buff
}
