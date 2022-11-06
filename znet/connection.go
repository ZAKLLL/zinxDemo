package znet

import (
	"fmt"
	io "io"
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
		Conn:         conn,
		ConnID:       connId,
		isClosed:     false,
		Router:       router,
		ExitBuffchan: make(chan bool),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	dp := &DataPack{}
	for {
		headData := make([]byte, dp.GetHeadLen())

		//buf := make([]byte, 512)
		_, err := io.ReadFull(c.Conn, headData)
		if err != nil {
			fmt.Println("recv Head err", err)
			c.ExitBuffchan <- true
			continue
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			c.ExitBuffchan <- true

			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*Message)
			msg.Data = make([]byte, msg.GetDataLen())
			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(c.Conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				c.ExitBuffchan <- true
				return
			}
			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen)

			req := Request{
				conn: c,
				msg:  msg,
			}
			go func(request ziface.IRequest) {
				c.Router.PreHandle(request)
				c.Router.Handle(request)
				c.Router.PostHandle(request)
			}(&req)
		}

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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	dp := &DataPack{}
	packData, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("db pack error ", err)
		c.ExitBuffchan <- true
		return err
	}
	if _, err := c.Conn.Write(packData); err != nil {
		fmt.Println("conn write error", err)
		c.ExitBuffchan <- true
		return err
	}
	return nil
}
