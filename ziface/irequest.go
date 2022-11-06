package ziface

type IRequest interface {
	GetConnection() IConnection //链接信息
	GetData() []byte
	GetMsgId() uint32
}
