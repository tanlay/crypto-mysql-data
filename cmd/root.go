package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

var RootCmd = &cobra.Command{
	Use:   "github.com/tanlay/crypto-mysql-data",
	Short: "解密mysql-data",
}
