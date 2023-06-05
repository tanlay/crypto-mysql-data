package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string //配置文件
	startId int64  //指定初始id
	limit   uint64 //限制查询条数
)

var RootCmd = &cobra.Command{
	Use:   "crypto-mysql-data",
	Short: "加解密mysql数据",
}
