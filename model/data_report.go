package model

type DataReport struct {
	Id             int64  `db:"id"`
	Cid            string `db:"cid"`
	CidType        int    `db:"cid_type"`
	Name           string `db:"name"`
	Phone          string `db:"phone"`
	DataSourceCode string `db:"data_source_code"`
	Data           string `db:"data"`
	Timestamp      int64  `db:"timestamp"`
	DataId         string `db:"data_id"`
	CreatedAt      int64  `db:"created_at"`
	DataGenTime    int64  `db:"data_gen_time"`
}

func (d *DataReport) TableName() string {
	return "data_report"
}

func (d *DataReport) DecryptedTableName() string {
	return "data_report_decrypted"
}
