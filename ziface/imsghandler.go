package ziface

type IMsgHandle interface {
	DoMsgHandle(request IRequest)           //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体处理逻辑
	StartWorkerPool()                       //启动worker工作池
	SendMsgToTaskQueue(request IRequest)    //发送待处理的任务给工作池
}
