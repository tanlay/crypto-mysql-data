package dao

import (
	"errors"
	"github.com/tanlay/crypto-mysql-data/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DataReportDecryptedDaoInterface interface {
	//PageById 按Id分页
	PageById(id int64, limit uint64) ([]model.DataReport, error)
	//Total 总数
	Total(id int64) (int, error)
	//BatchCreateDecrypted 解密
	BatchCreateDecrypted([]model.DataReport) error
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

// PageById 按id升序并分页排列
func (d *DataReportDecryptedDaoImpl) PageById(id int64, limit uint64) (reports []model.DataReport, err error) {
	dr := model.DataReport{}
	reports = make([]model.DataReport, 0)
	//统计data_report_decrypted表中，需要加密的数据，分页
	err = d.db.Table(dr.DecryptedTableName()).Select("*").Where("id > ?", id).
		Limit(int(limit)).Order("id asc").Find(&reports).Error
	if err != nil {
		d.logger.Error("page data_report err ", zap.Error(err))
		return nil, err
	}
	return reports, nil
}

// Total 查询需要加密的数据条数
func (d *DataReportDecryptedDaoImpl) Total(id int64) (int, error) {
	dr := model.DataReport{}
	var count int
	//统计data_report_decrypted表中，需要加密的数据量
	err := d.db.Select("count(*)").Table(dr.DecryptedTableName()).Where("id > ?", id).Find(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			d.logger.Error("get data_report count err", zap.Error(err))
			return 0, err
		}
	}
	return count, nil
}

// BatchCreateDecrypted 批量写入解密数据到data_report_decrypted表
func (d *DataReportDecryptedDaoImpl) BatchCreateDecrypted(dataReports []model.DataReport) error {
	dr := model.DataReport{}
	//写入解密数据到data_report_decrypted表
	err := d.db.Table(dr.DecryptedTableName()).Create(&dataReports).Error
	if err != nil {
		d.logger.Error("batch create data_report_decrypted err,", zap.Error(err))
		return err
	}
	return nil
}
