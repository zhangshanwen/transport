package rpc

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "github.com/zhangshanwen/transport/utils/proto"
)

func SendServer(homeRpc, addr, module string, pid int) (err error) {
	// Set up a connection to the server.
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(homeRpc, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewModuleClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Send(ctx, &pb.ModuleRequest{Addr: addr, Pid: int32(pid), Module: module})
	return
}

func Ping(addr string) (err error) {
	// Set up a connection to the server.
	var conn *grpc.ClientConn

	conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewModuleClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Ping(ctx, &pb.NormalRequest{Code: 1})
	return
}
