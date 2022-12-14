package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinxDemo/ziface"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinx.json来配置
*/
type GlobalObj struct {
	TcpServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称
	Version   string         //当前Zinx版本号

	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量

	MaxMsgChanLen int

	ConfigPath string
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile(g.ConfigPath)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}

}

/*
提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        7777,
		Host:           "0.0.0.0",
		MaxConn:        12000,
		MaxPacketSize:  4096,
		WorkerPoolSize: 1024,
		MaxMsgChanLen:  1024 * 10,
		ConfigPath:     "../conf/zinx.json",
	}

	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
