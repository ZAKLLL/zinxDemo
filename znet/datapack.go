package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinxDemo/utils"
	"zinxDemo/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	//id uInt32 四字节
	//dataLen uInt32 四字节
	return 8
}

//封包数据
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {

	dataBuff := bytes.NewBuffer([]byte{})

	//往databuff 里面写入datalen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写msgId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包数据
//行拆包的时候是分两次过程的，第二次是依赖第一次的dataLen结果，
//所以Unpack只能解压出包头head的内容，得到msgId 和 dataLen。
//之后调用者再根据dataLen继续从io流中读取body中的数据。
func (dp *DataPack) Unpack(binaryDate []byte) (ziface.IMessage, error) {
	dataBuffer := bytes.NewReader(binaryDate)

	msg := &Message{}

	// 读取dataLen 的值
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读取MessageId 的值
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen >= utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data recieved")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
