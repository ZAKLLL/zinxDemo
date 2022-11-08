package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"zinxDemo/utils"
	"zinxDemo/ziface"
)

/*
模拟客户端
*/
func ClientTest() {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	addr := fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	var idx uint32 = 0
	dp := &DataPack{}

	go func() {
		for {
			headData := make([]byte, dp.GetHeadLen())

			//buf := make([]byte, 512)
			_, err := io.ReadFull(conn, headData)
			if err != nil {
				fmt.Println("recv Head err", err)
				return
			}
			//将headData字节流 拆包到msg中
			msgHead, err := dp.Unpack(headData)
			if err != nil {
				fmt.Println("client unpack err:", err)
				return
			}

			if msgHead.GetDataLen() > 0 {
				msg := msgHead.(*Message)
				msg.Data = make([]byte, msg.GetDataLen())
				//根据dataLen从io中读取字节流
				_, err := io.ReadFull(conn, msg.Data)
				if err != nil {
					fmt.Println("client unpack data err:", err)
					return
				}
				fmt.Println("client===========> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))

			}

		}
	}()

	//发送数据
	go func() {
		for {
			msg := Message{
				Id:      0,
				DataLen: 0,
				Data:    nil,
			}
			idx++

			data := fmt.Sprintf("hello 现在是北京时间 %s", time.Now().Local().Format("2006/01/02 15:04:05"))
			bytes := []byte(data)
			msg.SetData(bytes)
			msg.SetDataLen(uint32(len(bytes)))
			packedData, err := dp.Pack(&msg)
			if err != nil {
				fmt.Println("pack msg error", err)
				return
			}

			_, err = conn.Write(packedData)

			if err != nil {
				fmt.Println("conn write error", err)
				return
			}

			time.Sleep(5 * time.Second)
		}
	}()

}

//该Server的连接创建时Hook函数
func MyOnConnStart(conn ziface.IConnection) {
	fmt.Println("hi conn", conn.GetConnID())
}

//该Server的连接断开时的Hook函数
func MyOnConnStop(conn ziface.IConnection) {

}

// Server 模块的测试函数
func TestServer(t *testing.T) {

	/*
		服务端测试
	*/
	//1 创建一个server 句柄 s
	s := NewServer()
	s.SetOnConnStart(MyOnConnStart)

	s.AddRouter(0, &MyRouter1{})
	/*
		客户端测试
	*/
	go ClientTest()
	go ClientTest()

	//2 开启服务
	s.Serve()
}
