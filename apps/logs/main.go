package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	module  = "logs"
	homeRpc = flag.String("rpc", ":50001", "")
)

func initialize() (err error) {
	flag.Parse() // 设置参数
	return
}

func main() {
	var err error
	flag.Parse() // 设置参数
	var lis net.Listener
	if lis, err = net.Listen("tcp", ":0"); err != nil {
		return
	}
	s := grpc.NewServer()
	if err = s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}

}
