package main

import(
	"net"
	"time"
)

type Listener struct{
	listener net.Listener
	connectionQueue chan struct{}
}

type Connection struct{
	conn net.Conn
	conneection	chan struct{}
}
func NewListner(listener_ net.Listener,maxConn int)*Listener{
	return &Listener{
		listener:listener_,
		connectionQueue:make(chan struct{},maxConn),
	};
}

func (l *Listener)Accept()(net.Conn,error){
	l.connectionQueue <-struct{}{}
	conn,err:=l.listener.Accept()
	if err!=nil{
		<-l.connectionQueue
		return nil,err
	}

	return &Connection{
		conn,
		l.connectionQueue,
	},nil
}

func (l*Listener)Close()error{
	return l.listener.Close()
}
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

func (conn *Connection)Close()error{
	err:=conn.conn.Close()
	<-conn.conneection
	return err
}

func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Connection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *Connection) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}