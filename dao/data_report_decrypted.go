package dao

import (
	"github.com/tanlay/crypto-mysql-data/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DataReportDecryptedDaoInterface interface {
	BatchCreate([]model.DataReport) error
}

type DataReportDecryptedDaoImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

var NewDataReportDecryptedDaoImpl = func(db *gorm.DB) DataReportDecryptedDaoInterface {
	return &DataReportDecryptedDaoImpl{
		db:     db,
		logger: zap.L().Named("data_report_decrypted_dao"),
	}
}

func (d *DataReportDecryptedDaoImpl) BatchCreate(dataReports []model.DataReport) error {
	dr := model.DataReport{}
	err := d.db.Table(dr.DecryptedTableName()).Create(&dataReports).Error
	if err != nil {
		d.logger.Error("batch create data_report err,", zap.Error(err))
		return err
	}
	return nil
}
