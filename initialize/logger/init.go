package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/zhangshanwen/transport/common"
)

var Writer io.Writer

func InitGinLogger() {
	gin.DefaultWriter = Writer
	logrus.Info("......GIN日志初始化成功......")
}

func InitLog(project, mod string, replica bool, index int) {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	filename := fmt.Sprintf("logs/%s.logs", project)
	if replica {
		filename = fmt.Sprintf("logs/%s_replca_%v.logs", project, index)
	}
	f := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1024, // megabytes
		MaxBackups: 10,
		MaxAge:     7, // days
	}
	if mod == common.ReleaseMode {
		Writer = io.MultiWriter(f)
	} else {
		Writer = io.MultiWriter(os.Stdout, f)
	}
	logrus.SetOutput(Writer)
}
