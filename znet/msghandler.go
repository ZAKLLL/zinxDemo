package znet

import (
	"fmt"
	"strconv"
	"zinxDemo/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法的map属性
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		//存放每个MsgId 所对应的处理方法的map属性
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	for !ok {
		fmt.Println("api msgId= ", request.GetMsgId(), "is not Found")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//不允许重复添加
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api,msgId=" + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId=", msgId)
}
