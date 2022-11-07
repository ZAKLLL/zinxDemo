package znet

import (
	"fmt"
	"net"
	"time"
	"zinxDemo/utils"
	"zinxDemo/ziface"
)

type Server struct {
	//服务器名称
	Name string
	//tcp4 or other
	IPVersion string
	// 服务绑定的IP地址
	IP string
	// 服务绑定的端口
	Port int
	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	//Router ziface.IRouter
	msgHandler ziface.IMsgHandle
}

type MyRouter1 struct {
	BaseRouter
}

//func (m MyRouter1) Handle(req ziface.IRequest) {
//
//	//回显业务
//	fmt.Println("server ------------------>" + string(req.GetData()))
//
//	backData := []byte("hi from zinxServer" + string(req.GetData()))
//
//	req.GetConnection().SendMsg(1, backData)
//}

func (s *Server) Start() {
	fmt.Printf("[START] Serve listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	go func() {

		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("Start Zinx server  ", s.Name, " succ, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32 = 0

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			dealConn := NewConnection(conn, cid, s.msgHandler)
			cid++
			//启动处理任务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//TODO  Serve.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

func (s *Server) Serve() {
	s.Start()

	//阻塞go主线程
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router for ", msgId, " succ! ")
}

func NewServer() ziface.IServer {
	config := utils.GlobalObject
	//config.Reload()

	s := &Server{
		Name:       config.Name,
		IPVersion:  "tcp4",
		IP:         config.Host,
		Port:       config.TcpPort,
		msgHandler: NewMsgHandle(),
	}
	config.TcpServer = s
	return s
}
