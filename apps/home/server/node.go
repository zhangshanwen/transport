package server

import (
	"fmt"
	"time"

	"github.com/zhangshanwen/transport/utils/rpc"
)

type Node struct {
	Ip      string `json:"ip"`
	Port    string `json:"port"`
	Rpc     string `json:"rpc"`
	Status  bool   `json:"status"`  // 运行状态 true 正常 0 停止
	Comment string `json:"comment"` // 备注
}

func (n *Node) getRpc() string {
	return fmt.Sprintf("%s:%v", n.Ip, n.Rpc)
}
func (n *Node) getAddr() string {
	return fmt.Sprintf("%s:%v", n.Ip, n.Port)
}

// ping 定时心跳检测
func (t *Transponder) ping() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			var err error
			for _, node := range t.nodes {
				// 检测节点是否启动
				// 检测节点是否已经为主节点，是则添加失败
				// 如果节点状态发生变化。分配规则发生变化
				node.Status = false
				addr := fmt.Sprintf("%s:%s", node.Ip, node.Port)
				if err = rpc.Ping(addr); err != nil {
					break
				}
				//var isMaster bool
				//if isMaster, err = rpc.Master(addr); err != nil {
				//	break
				//}
				//if isMaster {
				//	node.Comment = "该节点为主节点，无法完成添加"
				//}
				node.Status = true
				node.Comment = ""
			}
			<-ticker.C
		}
	}()
}

// 获取当前可用
func (t *Transponder) availableNodes() (nodes []Node) {
	for _, node := range t.nodes {
		if node.Status {
			nodes = append(nodes, *node)
		}
	}
	return
}

// dispense 分发模块
func (t *Transponder) dispense() (err error) {
	// 获取总计需要开启多少服务
	// 获取总计有多少节点
	// 服务分配规则
	/*
		1.如果只有主节点单台服务,则全部分配到主机节点
		2.如果节点有多点，则服务平均分配到除主节点以外的所有节点,主节点只做分发器
	*/
	nodes := t.availableNodes()
	if len(nodes) == 0 {
		// 当前无节点,本机执行任务
	} else {
		// 当前有节点,其他节点执行任务
	}
	return err
}
func (t *Transponder) syncFiles() {
	// 同步文件
	for _, node := range t.availableNodes() {
		// 判断主节点文件hash与节点hash是否相同
		for _, m := range t.Modules {
			m.Hash = ""
		}
		node.Comment = ""
	}
}
