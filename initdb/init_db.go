package initdb

import (
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/wrapper/store"
)

func InitDBDir(cfg *conf.Cfg) {
	InitSysSpace(cfg)
}

func InitSysSpace(conf *conf.Cfg) {
	store.NewSysTableSpace(conf)
}
