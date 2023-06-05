package controller

import (
	"errors"
	"fmt"
	"github.com/tanlay/crypto-mysql-data/config"
	"github.com/tanlay/crypto-mysql-data/dao"
	"github.com/tanlay/crypto-mysql-data/model"
	"github.com/tanlay/crypto-mysql-data/pkg/crypto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DataReportControllerInterface interface {
	DataReportBatchDecrypt(id int64, limit uint64) error
}

type DataReportControllerImpl struct {
	logger *zap.Logger
	db     *gorm.DB
}

var NewDataReportControllerImpl = func(db *gorm.DB) DataReportControllerInterface {
	return &DataReportControllerImpl{
		db:     db,
		logger: zap.L().Named("data_report_controller"),
	}
}

func (d *DataReportControllerImpl) DataReportBatchDecrypt(id int64, limit uint64) error {
	drDao := dao.NewDataReportDaoImpl(d.db)
	drdDao := dao.NewDataReportDecryptedDaoImpl(d.db)
	if limit > 1000 {
		return errors.New("超过单词查询的解密数量")
	}
	d.logger.Info("查询data_report需要解密的数量")
	count := 0

	total, err := drDao.Total(id)
	if err != nil {
		d.logger.Error("查询data_report数量错误", zap.Error(err))
		return err
	}
	if total == 0 {
		d.logger.Info("查询data_report没有需要解密的数据")
		return nil
	}
	d.logger.Info(fmt.Sprintf("查询data_report数量%d", total))
	d.logger.Info("开始解密")

	for {
		var decryptReports []model.DataReport
		reports, err := drDao.PageById(id, limit)
		if err != nil {
			d.logger.Error("查询data_report列表错误", zap.Error(err))
			return err
		}
		if len(reports) == 0 {
			break
		}
		for _, report := range reports {
			var err error
			dbSecretKey := config.C.Secret.DbSecretKey
			report.Cid, err = crypto.SM4ECBDecrypt(dbSecretKey, report.Cid)
			if err != nil {
				d.logger.Error("解密cid失败，", zap.Error(err), zap.Int64("id", report.Id))
				return err
			}
			report.Name, err = crypto.SM4ECBDecrypt(dbSecretKey, report.Name)
			if err != nil {
				d.logger.Error("解密name失败，", zap.Error(err), zap.Int64("id", report.Id))
				return err
			}
			report.Phone, err = crypto.SM4ECBDecrypt(dbSecretKey, report.Phone)
			if err != nil {
				d.logger.Error("解密phone失败，", zap.Error(err), zap.Int64("id", report.Id))
				return err
			}
			decryptReports = append(decryptReports, report)
			count++
		}
		if err := drdDao.BatchCreate(decryptReports); err != nil {
			d.logger.Error("写入解密数据错误", zap.Error(err))
			return err
		}
		id = reports[len(reports)-1].Id

		d.logger.Info(fmt.Sprintf("已解密%d,总数量%d,当前进度%.2f%%,最新一条数据id: %d",
			count, total, float64(count)/float64(total)*100, id))
	}
	d.logger.Info("数据解密完成")
	return nil
}
