package server

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"

	"github.com/zhangshanwen/transport/apps/home/conf"
	"github.com/zhangshanwen/transport/common"
	"github.com/zhangshanwen/transport/initialize/config"
	"github.com/zhangshanwen/transport/initialize/logger"
	"github.com/zhangshanwen/transport/utils"
	pb "github.com/zhangshanwen/transport/utils/proto"
	"github.com/zhangshanwen/transport/utils/rpc"
)

type (
	Transponder struct {
		mu                           sync.Mutex        // 读写锁
		Modules                      []*Module         // 模块
		pb.UnimplementedModuleServer                   // rpc 服务
		watcher                      *fsnotify.Watcher // 文件监控
		maxConnect                   int32             // 最大连接数
		connect                      int32             // 当前连接数
		connectTime                  int               // 连接时间
		lastConnectTime              *time.Time        // 最近连接时间
		nodes                        []*Node           // 所有节点
	}
)

/*
问:请求如何下发?
答:遍历模块,找到对应的prefix,然后根据index 找到不同的slave进行访问
问:怎么把模块的不同slave平均分配到各个节点?
答:为了保障主节点的性能,如果存在子节点，会有各个子节点平均分配各个模块的slave，如果slave的数量不平均，则前面的节点会承担更多的slave
问:怎么保障新增节点，或者节点挂掉，slave的分配呢?
答:如果新增节点,则根据算法优先从承担跟多salve的节点里面接手slave,如果节点挂掉，则把该节点挂载的slave平均分配到各个节点，优先从第一低的节点分配
问:如果各个节点都挂掉了呢?
答:将所有slave都挂载主节点上，如果有新的节点上线，则转移slave到该节点上
问:如何保障节点数发生变化时，各个模块的稳定性?
答:除非节点挂掉，会有短暂的无法访问,slave需要重新分配到各个节点进行重启,其余都会进行热启动，即保障服务的稳定的前提下进行重启
问:如何热启动?
答:重启启动一个slave关掉第一个slave,通过grpc通知到当前节点，当前节点然后通知主节点
问:如果按slave的数量进行平均分配,这样会不会更麻烦?
答:每个模块都会有一个优先级的，如果优先级越大，所占用的数量就越大，即能会被分配不同节点的几率就越大，默认优先级会被置为1，
比如有三个节点,有两个模块a,b，都有各有两个slave，a的优先级是2,b是默认，则a的两个slave会各分到一个节点，b的两个slave就只能分配到一个节点
问:这样分配有什么好处呢?
答:好处就是会让业务复杂且需要更多性能的节点更好的运行，而不会因为服务器的性能被其他应用占用而减缓处理速度
问:如果同时有大量的任务请求过来，服务会如何应对?
答:目前因为golang的net/http的瓶颈,会对请求进行限流，后续会持续优化
*/

func NewTransponder() (t *Transponder) {
	watcher, _ := fsnotify.NewWatcher()
	return &Transponder{
		watcher: watcher,
	}
}

func (t *Transponder) Run(project, mod string) (err error) {
	config.InitConf(project, &conf.C)      // 初始化配置文件
	logger.InitLog(project, mod, false, 0) // 初始化日志
	if err = copier.Copy(&t.Modules, &conf.C.Module); err != nil {
		return
	}
	//if err = copier.Copy(&server.Nodes, &conf.C.Nodes); err != nil {
	//	logrus.Error(err)
	//	return
	//}
	// TODO 节点管理
	//server.AddNodes() // 添加节点

	if len(t.Modules) == 0 {
		logrus.Warn("加载模块 0 个")
	}
	if err = config.WriteInConfig("pid", os.Getpid()); err != nil {
		logrus.Error(err)
		return
	}
	defer config.WriteInConfig("pid", 0)
	homeRpc := fmt.Sprintf("%s:%s", conf.C.Host, conf.C.Rpc)
	if err = t.RpcServer(homeRpc); err != nil {
		logrus.Error(err)
		return
	}
	if err = t.Listen(); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("转发 服务开启成功.......")

	for {
		// 等待gprc 服务完成开启
		if err = rpc.Ping(homeRpc); err == nil {
			logrus.Info("grpc 服务开启成功.......")
			break
		}
		logrus.Info("等待grpc服务开启.......")
		<-time.Tick(time.Millisecond * 500)
	}
	// 配置设置
	t.set(conf.C.MaxConnect, conf.C.ConnectTime)
	// 开启子程序
	t.StartUp(All, homeRpc)
	defer t.Close()
	return
}

// Close 关闭服务
func (t *Transponder) Close() {
	c := make(chan os.Signal)
	//监听所有信号
	signal.Notify(c, os.Kill, os.Interrupt)
	//阻塞直到有信号传入
	<-c
	_ = t.watcher.Close()
	t.ShutDown(All)
	logrus.Info("主程序完成退出......")
}

// 设置最大连接数，连接时长
func (t *Transponder) set(maxConnect int32, connectTime int) {
	if maxConnect <= 0 {
		maxConnect = common.MaxConnect
	}
	t.maxConnect = maxConnect
	if connectTime <= 0 {
		connectTime = common.ConnectTime
	}
	t.connectTime = connectTime
}

// 检查文件并设置文件hash值
func (t *Transponder) checkFiles() (err error) {
	for _, m := range t.Modules {
		if m.Hash, err = utils.GetFileHash(m.Cmd); err != nil {
			return
		}
	}
	return
}
