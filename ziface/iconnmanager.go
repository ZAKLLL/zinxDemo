package ziface

type IConnManager interface {
	Add(conn IConnection)                   //添加链接
	Remove(conn IConnection)                //删除链接
	Get(connId uint32) (IConnection, error) //获取链接
	Len() int                               //获取当前链接数量
	ClearConn()                             //删除并停止所有链接
	IsFull() bool                           //连接池是否满了
}
