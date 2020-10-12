package xenv

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"j2pay-server/hcommon"
	"j2pay-server/pkg/setting"
)

// DbCon 数据库链接
var DbCon *sqlx.DB

func EnvCreate() {
	// 初始化数据库
	sql := fmt.Sprintf("%s",setting.SqlxConf.Name)
	DbCon = hcommon.DbCreate(sql, true)
}
// EnvDestroy 销毁运行环境
func EnvDestroy() {
	if DbCon != nil {
		_ = DbCon.Close()
	}
}
