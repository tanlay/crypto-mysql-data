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
	DataReportBatchEncrypt(id int64, limit uint64) error
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
	dataReportDao := dao.NewDataReportDaoImpl(d.db)
	dataReportDecryptedDao := dao.NewDataReportDecryptedDaoImpl(d.db)
	if limit > 1000 {
		return errors.New("超过单词查询的解密数量")
	}
	d.logger.Info("查询data_report需要解密的数量")
	decryptCount := 0

	total, err := dataReportDao.Total(id)
	if err != nil {
		d.logger.Error("查询data_report数量错误", zap.Error(err))
		return err
	}
	if total == 0 {
		d.logger.Info("查询data_report没有需要解密的数据")
		return nil
	}
	d.logger.Info(fmt.Sprintf("查询data_report数量%d", total))
	d.logger.Info("data_report开始解密")

	for {
		var decryptReports []model.DataReport
		reports, err := dataReportDao.PageById(id, limit)
		if err != nil {
			d.logger.Error("查询data_report列表错误", zap.Error(err))
			return err
		}
		if len(reports) == 0 {
			break
		}
		for _, report := range reports {
			var err error
			dbSecretKey := config.C.Secret.DecSecretKey //解密秘钥
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
			decryptCount++
		}
		if err := dataReportDecryptedDao.BatchCreateDecrypted(decryptReports); err != nil {
			d.logger.Error("写入data_report_decrypted解密数据错误", zap.Error(err))
			return err
		}
		id = reports[len(reports)-1].Id

		d.logger.Info(fmt.Sprintf("已解密%d,总数量%d,当前进度%.2f%%,最新一条数据id: %d",
			decryptCount, total, float64(decryptCount)/float64(total)*100, id))
	}
	d.logger.Info("数据解密完成")
	return nil
}

func (d *DataReportControllerImpl) DataReportBatchEncrypt(id int64, limit uint64) error {
	dataReportDecrypted := dao.NewDataReportDaoImpl(d.db)
	dataReportEncrypted := dao.NewDataReportEncryptedDaoImpl(d.db)
	if limit > 1000 {
		return errors.New("超过单词查询的加密数量")
	}
	d.logger.Info("查询data_report_decrypted需要加密的数量")
	encryptCount := 0

	total, err := dataReportDecrypted.Total(id)
	if err != nil {
		d.logger.Error("查询data_report_decrypted数量错误", zap.Error(err))
		return err
	}
	if total == 0 {
		d.logger.Info("查询data_report_decrypted没有需要加密的数据")
		return nil
	}
	d.logger.Info(fmt.Sprintf("查询data_report_decrypted数量%d", total))
	d.logger.Info("开始加密data_report_decrypted")

	for {
		var encryptReports []model.DataReport
		reports, err := dataReportDecrypted.PageById(id, limit)
		if err != nil {
			d.logger.Error("查询data_report_decrypted列表错误", zap.Error(err))
			return err
		}
		if len(reports) == 0 {
			break
		}
		for _, report := range reports {
			var err error
			dbSecretKey := config.C.Secret.EncSecretKey //加密秘钥
			report.Cid, err = crypto.SM4ECBEncrypt(dbSecretKey, report.Cid)
			if err != nil {
				d.logger.Error("解密cid失败，", zap.Error(err), zap.Int64("id", report.Id))
				return err
			}
			report.Name, err = crypto.SM4ECBEncrypt(dbSecretKey, report.Name)
			if err != nil {
				d.logger.Error("解密name失败，", zap.Error(err), zap.Int64("id", report.Id))
				return err
			}
			report.Phone, err = crypto.SM4ECBEncrypt(dbSecretKey, report.Phone)
			if err != nil {
				d.logger.Error("解密phone失败，", zap.Error(err), zap.Int64("id", report.Id))
				return err
			}
			encryptReports = append(encryptReports, report)
			encryptCount++
		}
		if err := dataReportEncrypted.BatchCreateEncrypted(encryptReports); err != nil {
			d.logger.Error("写入data_report_encrypted加密数据错误", zap.Error(err))
			return err
		}
		id = reports[len(reports)-1].Id

		d.logger.Info(fmt.Sprintf("已加密%d,总数量%d,当前进度%.2f%%,最新一条数据id: %d",
			encryptCount, total, float64(encryptCount)/float64(total)*100, id))
	}
	d.logger.Info("数据加密完成")
	return nil
}
