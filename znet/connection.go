package znet

import (
	"fmt"
	"net"
	"zinxDemo/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool

	//handleApi ziface.HandFunc

	//该连接的处理方法router
	Router ziface.IRouter

	//告知该链接已经退出/停止的channel
	ExitBuffchan chan bool
}

//创建新链接
func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connId,
		isClosed: false,
		//handleApi:    callback_api,
		Router:       router,
		ExitBuffchan: make(chan bool),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			c.ExitBuffchan <- true
			continue
		}
		req := Request{
			conn: c,
			data: buf,
		}
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		//if err := c.handleApi(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnId", c.ConnID, " handle is error")
		//	c.ExitBuffchan <- true
		//	return
		//}
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	for {
		select {
		case <-c.ExitBuffchan:
			//接收到退出的消息直接退出
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	c.Conn.Close()
	c.ExitBuffchan <- true
	close(c.ExitBuffchan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
