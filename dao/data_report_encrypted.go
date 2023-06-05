package dao

import (
	"github.com/tanlay/crypto-mysql-data/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DataReportEncryptedDaoInterface interface {
	//BatchCreateEncrypted 加密
	BatchCreateEncrypted([]model.DataReport) error
}

type DataReportEncryptedDaoImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

var NewDataReportEncryptedDaoImpl = func(db *gorm.DB) DataReportEncryptedDaoInterface {
	return &DataReportEncryptedDaoImpl{
		db:     db,
		logger: zap.L().Named("data_report_encrypted_dao"),
	}
}

// BatchCreateEncrypted 批量写入加密数据到data_report_encrypted表
func (d *DataReportEncryptedDaoImpl) BatchCreateEncrypted(dataReports []model.DataReport) error {
	dr := model.DataReport{}
	//写入解密数据到data_report_encrypted表
	err := d.db.Table(dr.EncryptedTableName()).Create(&dataReports).Error
	if err != nil {
		d.logger.Error("batch create data_report_encrypted err,", zap.Error(err))
		return err
	}
	return nil
}
