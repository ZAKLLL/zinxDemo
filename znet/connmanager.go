package znet

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"zinxDemo/utils"
	"zinxDemo/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex //读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", cm.Len())
}
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//删除连接信息
	delete(cm.connections, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", cm.Len())

}
func (cm *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found by connId: " + strconv.Itoa(int(connId)))
}
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connId, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connId)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) IsFull() bool {
	return cm.Len() >= utils.GlobalObject.MaxConn
}
