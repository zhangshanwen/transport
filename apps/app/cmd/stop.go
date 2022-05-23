/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhangshanwen/transport/common"
	"github.com/zhangshanwen/transport/utils"
	"io/ioutil"
	"strconv"
	"syscall"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止",
	Run: func(cmd *cobra.Command, args []string) {
		// 根据 pid 关闭程序，并检测程序是否关闭
		var content []byte
		var err error
		if content, err = ioutil.ReadFile(common.HomePid); err != nil {
			utils.Println("red", "关闭程序失败", err.Error())
			return
		}
		var pid int
		if pid, err = strconv.Atoi(string(content)); err != nil {
			utils.Println("red", "pid错误", err.Error())
			return
		}
		if err = syscall.Kill(pid, syscall.SIGINT); err != nil {
			utils.Println("red", "程序关闭失败", err.Error())
			return
		}
		utils.Println("green", "程序关闭成功........")
	},
}
