package main

import (
	"github.com/zhangshanwen/transport/apps/home/server"
)

var (
	project = "home"
	mod     = "debug"
)

func main() {
	t := server.NewTransponder()
	if err := t.Run(project, mod); err != nil {
		panic(err)
	}
}
