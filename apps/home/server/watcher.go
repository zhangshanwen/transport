package server

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"time"
)

func (t *Transponder) watching(module *Module, homeRpc string) {
	if err := t.watcher.Add(module.Cmd); err != nil {
		logrus.Errorf("模块%s,文件%s失败.........", module.Name, module.Cmd)
		return
	}
	logrus.Infof("模块%s,文件%s监视中.........", module.Name, module.Cmd)
	go func() {
		for {
			select {
			case event := <-t.watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					// 文件修改重启发送重启改服务指令
					t.StartUp(module.Name, homeRpc)
					logrus.Info("发送开启指令...............", event.Op.String())
					time.Sleep(1 * time.Second)
				}
			case err := <-t.watcher.Errors:
				logrus.Error("error:", err)
			}
		}
	}()

}
