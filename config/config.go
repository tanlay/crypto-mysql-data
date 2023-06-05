package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/tanlay/crypto-mysql-data/pkg/db"
	"go.uber.org/zap"
)

var C = &Config{}

type Config struct {
	Database db.DatabaseConf `json:"database" toml:"database"`
	Secret   SecretConf      `json:"secret" toml:"secret"`
	Logger   LoggerConf      `json:"logger" toml:"logger"`
}

type SecretConf struct {
	DbSecretKey string `json:"db_secret_key" toml:"db_secret_key"`
}

type LoggerConf struct {
	Env    string `json:"env" toml:"env"`
	Level  string `json:"level" toml:"level"`
	Output string `json:"output" toml:"output"`
}

func LoadCongFromToml(cfgFile string) error {
	if cfgFile == "" {
		return errors.New("需指定配置文件")
	} else {
		if _, err := toml.DecodeFile(cfgFile, C); err != nil {
			return err
		}
	}
	zapLevel := zap.NewAtomicLevel()
	if err := zapLevel.UnmarshalText([]byte(C.Logger.Level)); err != nil {
		panic(err.Error())
	}
	var zapConf zap.Config

	if env := C.Logger.Env; env == "dev" {
		zapConf = zap.NewDevelopmentConfig()
	} else {
		zapConf = zap.NewProductionConfig()
	}
	zapConf.Level = zapLevel

	if C.Logger.Output != "" {
		zapConf.OutputPaths = []string{C.Logger.Output}
		zapConf.ErrorOutputPaths = []string{C.Logger.Output}
	}
	if logger, err := zapConf.Build(); err != nil {
		panic(err.Error())
	} else {
		zap.RedirectStdLog(logger)
		zap.ReplaceGlobals(logger)
	}
	zap.L().Info(fmt.Sprintf("load config: %s", cfgFile))
	return nil
}
