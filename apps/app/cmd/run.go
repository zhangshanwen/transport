/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhangshanwen/transport/utils"
	"github.com/zhangshanwen/transport/utils/rpc"
	"os/exec"
)

/*
 */

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if err = rpc.Ping(); err == nil {
			utils.Println("red", "程序正在运行中........")
			return
		}
		// 开启服务
		// 判断pid 文件是否存在，存在则删除
		if err = exec.Command("nohup", "bin/home").Start(); err != nil {
			utils.Println("red", "程序运行失败", err.Error())
			return
		}
		utils.Println("green", "程序开启成功.....")
	},
}
