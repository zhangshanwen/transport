package server

import (
	"github.com/sirupsen/logrus"
)

const (
	All = "all"
)

type (
	Module struct {
		Name       string
		Replica    int
		Cmd        string
		Hash       string
		Prefix     string
		Scheme     string
		IsWatching bool
		Index      int32
		Slaves     []Slave
	}
	Slave struct {
		Addr string
		Pid  int
	}
)

func (t *Transponder) ShutDown(name string) {
	var err error
	for _, module := range t.Modules {
		if name != module.Name && name != All {
			continue
		}
		logrus.Infof("开始关闭%s,副本个数:%v ......", module.Name, module.Replica)
		for i, slave := range module.Slaves {
			logrus.Infof("开始关闭 %s,第%v个副本 ......", module.Name, i+1)
			if err = t.killCmd(slave.Pid); err != nil {
				logrus.Errorf("关闭失败%s,第%v个副本 err:%v ......", module.Name, i+1, err)
			} else {
				logrus.Infof("关闭成功%s,第%v个副本 ......", module.Name, i+1)
			}
		}
	}
}

func (t *Transponder) StartUp(name, homeRpc string) {
	logrus.Infof("开始加载模块%s......", name)
	var err error
	for _, module := range t.Modules {
		if name != module.Name && name != All {
			continue
		}
		// 开启监视文件指令
		if !module.IsWatching {
			t.watching(module, homeRpc)
			module.IsWatching = true
		}

		logrus.Infof("开始启动%s,副本个数:%v ......", module.Name, module.Replica)
		for i := 0; i < module.Replica+1; i++ {
			logrus.Infof("开始启动 %s,第%v个副本 ......", module.Name, i+1)
			if err = t.startCmd(module.Cmd, homeRpc, i+1, module.Replica); err != nil {
				logrus.Errorf("启动失败%s,第%v个副本 err:%v ......", module.Name, i+1, err)
				continue
			} else {
				logrus.Infof("启动成功%s,第%v个副本 ......", module.Name, i+1)
			}
			if len(module.Slaves) > 0 {
				// 热启动
				slave := module.Slaves[0] // 获取第一个节点
				module.Slaves = module.Slaves[1:]
				if err = t.killCmd(slave.Pid); err != nil {
					logrus.Errorf("热启动,关闭失败%s,第%v个副本 err:%v ......", module.Name, i+1, err)
				} else {
					logrus.Infof("热启动,关闭成功%s,第%v个副本 ......", module.Name, i+1)
				}

			}
		}
	}
	logrus.Infof("模块%s加载完成......", name)
}

// 获取服务个数
func (t *Transponder) getModuleLength() (l int) {
	for _, m := range t.Modules {
		l += m.Replica + 1
	}
	return
}
