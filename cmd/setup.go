package cmd

import (
	"github.com/tanlay/crypto-mysql-data/config"
	"github.com/tanlay/crypto-mysql-data/pkg/db"
)

func SetupGlobalDB(conf config.Config) error {
	return db.GetGlobalDB(conf.Database)
}
