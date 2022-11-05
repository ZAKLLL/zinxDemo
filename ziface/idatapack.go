package ziface

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error) //封包
	Unpack([]byte) (IMessage, error)   //拆包

}
