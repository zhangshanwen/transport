package rpc

import (
	"context"
	"time"

	pb "github.com/zhangshanwen/transport/utils/proto"

	"google.golang.org/grpc"
)

func SyncFiles(addr, filename string, b []byte) (err error) {
	// Set up a connection to the server.
	var conn *grpc.ClientConn

	conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.SyncFile(ctx, &pb.FileRequest{Name: filename, Files: b})
	return
}

func Hash(addr, filename string) (err error) {
	// Set up a connection to the server.
	var conn *grpc.ClientConn

	conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Hash(ctx, &pb.HashRequest{Name: filename})
	return
}
