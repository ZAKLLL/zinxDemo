package znet

import "zinxDemo/ziface"

type Request struct {
	conn ziface.IConnection //已经建立好的链接
	//data []byte             //客户端请求的数据
	msg ziface.IMessage //客户端请求的数据
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsg() uint32 {
	return r.msg.GetMsgId()
}
