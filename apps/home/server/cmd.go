package server

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

func (t *Transponder) killCmd(pid int) (err error) {
	if err = syscall.Kill(pid, syscall.SIGINT); err != nil {
		return
	}
	logrus.Infof("关闭进程成功 pid %d ", pid)
	return
}

func (t *Transponder) startCmd(cmd, homeRpc string, index, replica int) (err error) {
	c := exec.Command("nohup", cmd, fmt.Sprintf("-rpc=%s", homeRpc), fmt.Sprintf("-replica=%v", replica > 1), fmt.Sprintf("-index=%v", index))
	return c.Start()
}
