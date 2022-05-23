package server

import (
	"context"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/zhangshanwen/transport/utils/proto"
)

func (t *Transponder) Send(ctx context.Context, in *pb.ModuleRequest) (*pb.ModuleReply, error) {
	for _, module := range t.Modules {
		if module.Name == in.GetModule() {
			module.Slaves = append(module.Slaves, Slave{
				Addr: in.GetAddr(),
				Pid:  int(in.GetPid()),
			})
		}
	}
	logrus.Infof("模块:%s创建成功,地址:%s,pid:%v", in.GetModule(), in.GetAddr(), in.GetPid())
	return &pb.ModuleReply{Message: "ok"}, nil
}
func (t *Transponder) Ping(ctx context.Context, in *pb.NormalRequest) (*pb.NormalReply, error) {
	return &pb.NormalReply{Code: in.GetCode()}, nil
}
func (t *Transponder) Master(ctx context.Context, in *pb.NormalRequest) (*pb.NormalReply, error) {
	return &pb.NormalReply{Code: int32(len(t.nodes))}, nil
}

func (t *Transponder) RpcServer(addr string) (err error) {
	var lis net.Listener
	if lis, err = net.Listen("tcp", addr); err != nil {
		return
	}
	s := grpc.NewServer()
	pb.RegisterModuleServer(s, t)
	go func() {
		if err = s.Serve(lis); err != nil {
			log.Panicf("failed to serve: %v", err)
		}
	}()
	return
}
