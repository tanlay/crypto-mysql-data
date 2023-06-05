package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tanlay/crypto-mysql-data/config"
	"github.com/tanlay/crypto-mysql-data/controller"
	"github.com/tanlay/crypto-mysql-data/pkg/db"
	"go.uber.org/zap"
)

func init() {
	RootCmd.AddCommand(DecryptCmd())
}

func DecryptCmd() *cobra.Command {
	limit = 100 //设置默认limit为100
	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "解密",
		Long:  "解密",
		RunE: func(cmd *cobra.Command, args []string) error {
			tasks := []func(config config.Config) error{
				SetupGlobalDB,
			}
			for _, task := range tasks {
				if err := task(*config.C); err != nil {
					zap.L().Panic("setup failed", zap.Error(err))
					return err
				}
			}
			svr := controller.NewDataReportControllerImpl(db.GlobalDB)
			if err := svr.DataReportBatchDecrypt(startId, limit); err != nil {
				cmd.PrintErrf(err.Error())
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.LoadCongFromToml(cfgFile)
		},
	}
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "配置文件目录")
	cmd.PersistentFlags().Int64VarP(&startId, "start_id", "s", 0, "开始ID编号")
	cmd.PersistentFlags().Uint64VarP(&limit, "limit", "l", 100, "查询数量（1-1000）,默认100")
	return cmd
}
