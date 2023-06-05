package dao

import (
	"errors"
	"github.com/tanlay/crypto-mysql-data/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DataReportDaoInterface interface {
	//PageById 按Id分页
	PageById(id int64, limit uint64) ([]model.DataReport, error)
	//Total 总数
	Total(id int64) (int, error)
}

type DataReportDaoImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

var NewDataReportDaoImpl = func(db *gorm.DB) DataReportDaoInterface {
	return &DataReportDaoImpl{
		db:     db,
		logger: zap.L().Named("data_report_dao"),
	}
}

func (d *DataReportDaoImpl) PageById(id int64, limit uint64) (reports []model.DataReport, err error) {
	dr := model.DataReport{}
	reports = make([]model.DataReport, 0)
	err = d.db.Table(dr.TableName()).Select("*").Where("id > ?", id).
		Limit(int(limit)).Order("id asc").Find(&reports).Error
	if err != nil {
		d.logger.Error("page data_report err ", zap.Error(err))
		return nil, err
	}
	return reports, nil
}

func (d *DataReportDaoImpl) Total(id int64) (int, error) {
	dr := model.DataReport{}
	var count int
	err := d.db.Select("count(*)").Table(dr.TableName()).Where("id > ?", id).Find(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			d.logger.Error("get data_report count err", zap.Error(err))
			return 0, err
		}
	}
	return count, nil
}
