package main
import(
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"strconv"
)
//
type SocketConn struct {
	name  	string
	conn    net.Conn
	scanner *bufio.Scanner
	lock 	sync.RWMutex
}


func (c *SocketConn) Write(msg ...string) error {
	//提前上锁,防止出现意外
	c.lock.Lock()
	defer c.lock.Unlock()
	//
	if c.conn==nil{
		return errors.New("当前无连接")
	}
	//
	for _, arg := range msg {
		_, err := fmt.Fprintf(c.conn, "%s", arg+"\n")
		if err != nil {
			return errors.New("写入失败")
		}
	}
	return nil
}
func (c *SocketConn) Read() (string, error) {
	//读取不需要加互斥锁
	if c.scanner.Scan() {
		return c.scanner.Text(), nil
	}
	return "", c.scanner.Err()
}
func (c *SocketConn) SetConn(conn net.Conn,scanner *bufio.Scanner){
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn!=nil{
		c.conn.Close()
	}
	//
	c.conn=conn
	c.scanner=scanner
}
func MakeSocketConn(name string, port int) (*SocketConn, error) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	//接受客户端连接
	Debug("开始监听端口:"+strconv.Itoa(port))
	ret := &SocketConn{
			name:     name,
			conn:    nil,
			scanner: nil,
	}
	go func(){
		for{
			conn, err := listener.Accept()
			//
			if err != nil {
				Debug("连接异常")
				continue
			}
			//如果能读取到通用cookie，则更换连接
			scanner:=bufio.NewScanner(conn)
			if scanner.Scan(){
				if scanner.Text()==UniverseCookie{
					ret.SetConn(conn,scanner)
					Debug("端口"+strconv.Itoa(port)+"连接成功!")
				}
			}else{
				conn.Close()
			}
		}
	}()
	//
	return ret, nil
}