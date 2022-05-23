package main

import (
	"flag"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"

	"github.com/zhangshanwen/transport/apps/admin/conf"
	"github.com/zhangshanwen/transport/apps/admin/migrate"
	"github.com/zhangshanwen/transport/apps/admin/router"
	"github.com/zhangshanwen/transport/apps/admin/tools"
	"github.com/zhangshanwen/transport/initialize/app"
	"github.com/zhangshanwen/transport/initialize/config"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/logger"
	"github.com/zhangshanwen/transport/initialize/node"
	"github.com/zhangshanwen/transport/utils/rpc"
)

var (
	module  = "admin"
	mod     = "debug"
	replica = flag.Bool("replica", false, "")
	index   = flag.Int("index", 0, "")
	homeRpc = flag.String("rpc", ":50001", "")
)

func initialize() (err error) {
	gin.SetMode(mod)                              // 设置环境
	flag.Parse()                                  // 设置参数
	logger.InitLog(module, mod, *replica, *index) // 初始化日志
	config.InitConf(module, &conf.C)              // 初始化配置
	tools.InitJwt(module)                         // 初始jwt
	mysql := db.Mysql{}
	if err = copier.Copy(&mysql, &conf.C.DB.Mysql); err != nil {
		logrus.Error("复制失败", err)
		return
	}
	db.InitMysql(mysql) // 初始化mysql
	redis := db.Redis{}
	if err = copier.Copy(&redis, &conf.C.DB.Redis); err != nil {
		logrus.Error("复制失败", err)
		return
	}
	db.InitRedis(redis)     // 初始redis
	node.InitNode()         // 初始node
	router.RegisterRouter() // 注册路由
	app.InitRoute(module)   // 初始化路由
	migrate.AutoMigrate()
	return
}

func main() {
	var err error
	if err = initialize(); err != nil {
		logrus.Error("初始化失败", err)
		return
	}
	var l net.Listener
	if l, err = net.Listen("tcp", ":0"); err != nil {
		logrus.Error("创建监听端口失败", err)
		return
	}
	addr := l.Addr().String()
	logrus.Info("监听http服务 ", addr)
	pid := os.Getpid()
	go func() {
		if err = rpc.SendServer(*homeRpc, addr, module, pid); err != nil {
			logrus.Error("向主服务发送基本信息失败 ", err)
		} // 向主服务发送基本信息
	}()
	if err = http.Serve(l, app.R); err != nil {
		logrus.Error("服务开启失败 ", err)
		return
	}

}
