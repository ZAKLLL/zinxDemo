package ziface

import "net"

//连接接口定义
type IConnection interface {
	Start()
	Stop()
	GetTcpConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
}

//统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
