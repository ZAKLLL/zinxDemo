package znet

import (
	"errors"
	"fmt"
	io "io"
	"net"
	"zinxDemo/utils"
	"zinxDemo/ziface"
)

type Connection struct {
	//隶属于哪个服务
	TcpServer ziface.IServer

	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool

	//handleApi ziface.HandFunc

	MsgHandler ziface.IMsgHandle

	//告知该链接已经退出/停止的channel
	ExitBuffchan chan bool

	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte

	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
}

//创建新链接
func NewConnection(server ziface.IServer, conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connId,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffchan: make(chan bool),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen), //不要忘记初始化

	}
	//将当前链接关联到server的链接管理器中
	server.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	dp := &DataPack{}
	for {

		headData := make([]byte, dp.GetHeadLen())

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
			if utils.GlobalObject.WorkerPoolSize > 0 {
				c.MsgHandler.SendMsgToTaskQueue(&req)
			}
		}

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer GoRoutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				c.ExitBuffchan <- true
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffchan:
			//conn 已经关闭
			return
		}
	}

}

func (c *Connection) Start() {
	//启动读
	go c.StartReader()

	go c.StartWriter()

	c.TcpServer.CallOnConnStart(c)

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
	c.TcpServer.CallOnConnStop(c)

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
	//将消息发送到writechannel 中让Writer发送给客户端
	c.msgChan <- packData

	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	dp := &DataPack{}
	packData, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("db pack error ", err)
		c.ExitBuffchan <- true
		return err
	}
	c.msgBuffChan <- packData

	return nil
}
