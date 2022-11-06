package datapackTest_test

import (
	"fmt"
	"net"
	"testing"
	"time"
	"zinxDemo/znet"
)

func TestClient(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("net dial error", err)
	}

	tmpSendData := make([]byte, 0)

	var idx uint32 = 0
	dp := &znet.DataPack{}
	for {
		msg := znet.Message{
			Id:      idx,
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
		tmpSendData = append(tmpSendData, packedData...)

		if idx%2 == 0 {

			_, err = conn.Write(tmpSendData)

			if err != nil {
				fmt.Println("conn write error", err)
				return
			}
			tmpSendData = make([]byte, 0)

		}

		time.Sleep(2 * time.Second)
	}

}
