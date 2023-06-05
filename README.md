# crypto-mysql-data

## 功能
实现数据表某些字段的加解密并存储到另一张表中
- -c "config.toml"，指定配置文件
- -s "id",来继续执行上一次结束的任务
- -l 1000,来指定单次查询的数量，取值为0-1000，默认100

## 安装依赖
```shell
go mod tidy
go build
```

```shell
./crypto-mysql-data decrypt -h
解密

Usage:
  github.com/tanlay/crypto-mysql-data decrypt [flags]

Flags:
  -c, --config string   配置文件目录 (default "config.toml")
  -h, --help            help for decrypt
  -l, --limit uint      查询数量(0-1000),默认100 (default 100)
  -s, --start_id int    开始ID编号
```

## 解密id大于2527503360441305664的所有记录
```go
./crypto-mysql-data decrypt -c config.toml -l 1000 -s 2527503360441305664
```

## conf.toml配置文件
```toml
[database]
dsn="mysql://root:123456@tcp(localhost:3306)/yt_judgery_prod?parseTime=True"

[secret]
db_secret_key="b44982c8d69333e976e13cbb3ba127fb"

[logger]
env="prod"
level="debug"
output="log.txt"
```
