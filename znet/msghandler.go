package znet

import (
	"fmt"
	"strconv"
	"zinxDemo/utils"
	"zinxDemo/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法的map属性

	WorkerPoolSize uint32 //业务工作worker 池的数量

	TaskQueue []chan ziface.IRequest //Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		//存放每个MsgId 所对应的处理方法的map属性
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	for !ok {
		fmt.Println("api msgId= ", request.GetMsgId(), "is not Found , closing this ")
		request.GetConnection().Stop()
		return
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

func (mh *MsgHandle) startOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerId, " is started.")
	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandle(req)
		}
	}

}

//启动工作池
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i += 1 {
		//开启一个worker
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerId := request.GetConnection().GetConnID() % utils.GlobalObject.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgId(), "to workerID=", workerId)
	mh.TaskQueue[workerId] <- request
}
