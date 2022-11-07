package ziface

import "net"

//连接接口定义
type IConnection interface {
	Start()
	Stop()
	GetTcpConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error
}
