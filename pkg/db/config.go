package db

type DatabaseConf struct {
	DSN             string `json:"dsn"`
	MaxIdleConn     int    `json:"max_idle_conn"`
	MaxOpenConn     int    `json:"max_open_conn"`
	ConnMaxLiftTime int    `json:"conn_max_lift_time"`
}
