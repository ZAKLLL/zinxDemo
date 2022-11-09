package ziface

import (
	"net"
)

// 连接接口定义
type IConnection interface {
	Start() //开始处理链接
	Stop()  //停止处理链接
	GetTcpConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error     //发送消息，无缓冲
	SendBuffMsg(msgId uint32, data []byte) error //发送缓冲消息

	//设置属性值
	SetProperty(key string, value interface{})

	//读取属性值
	GetProperty(key string) (interface{}, error)

	//删除属性值
	RemoveProperty(key string)
}
