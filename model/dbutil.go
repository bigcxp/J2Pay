package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"j2pay-server/pkg/setting"
	"log"
)

var db *gorm.DB

type Base struct{}
type FieldTrans map[string]string
func Setup() {
	fmt.Println("init db")
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.MysqlConf.User,
		setting.MysqlConf.Pwd,
		setting.MysqlConf.Host,
		setting.MysqlConf.Port,
		setting.MysqlConf.Db))
	if err != nil {
		log.Panicf("连接数据库错误 ：%s", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.MysqlConf.Prefix + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.SetLogger(&GormLogger{})
	db.DB().SetMaxIdleConns(setting.MysqlConf.MaxIdle)
	db.DB().SetMaxOpenConns(setting.MysqlConf.MaxActive)
	AutoMigrate()
	// 设置程序启动参数 -init | -init=true
	if setting.Init {
		InitSql()
	}
}


func Getdb() *gorm.DB{
	if db == nil{
		Setup()
	}
	return db
}

// 通用分页获取偏移量
func GetOffset(page, pageSize int) int {
	if page <= 1 {
		return 0
	}
	return (page - 1) * pageSize
}


// 设置条件
func MultiWhere(where ...interface{}) func(db2 *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(where[0], where[1:]...)
	}
}


// 设置条件
func MultiOr(where ...interface{}) func(db2 *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(where[0], where[1:]...)
	}
}


// 自动创建修改表
func AutoMigrate() {
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '后台用户'").AutoMigrate(&AdminUser{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '角色'").AutoMigrate(&Role{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '权限'").AutoMigrate(&Auth{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT 'casbin policy 配置'").AutoMigrate(&CasbinRule{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '系统公告'").AutoMigrate(&SystemMessage{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '商户提领'").AutoMigrate(&Pick{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '订单'").AutoMigrate(&Order{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '退款订单'").AutoMigrate(&Return{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '手续费结账'").AutoMigrate(&Fee{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '实收明细记录'").AutoMigrate(&DetailedRecord{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '系统参数设定'").AutoMigrate(&Parameter{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '汇率表'").AutoMigrate(&Rate{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '收款地址表'").AutoMigrate(&Address{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '默认整型配置表'").AutoMigrate(&TAppConfigInt{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '程序状态表'").AutoMigrate(&TAppStatusInt{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '默认字符串配置表'").AutoMigrate(&TAppConfigStr{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT 'eth充币交易表'").AutoMigrate(&TTx{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT 'eRC20交易表'").AutoMigrate(&TTxErc20{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '待发送表'").AutoMigrate(&TSend{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '提币记录表'").AutoMigrate(&TWithdraw{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '通知表'").AutoMigrate(&TProductNotify{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '合约配置表'").AutoMigrate(&TAppConfigToken{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT 'eth交易记录表'").AutoMigrate(&EthTransaction{})
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8 COMMENT '热钱包交易记录表'").AutoMigrate(&HotTransaction{})



}

func InitSql() {
	// 清空
	//db.Exec("truncate admin_user")
	//db.Exec("truncate role")
	//db.Exec("truncate casbin_rule")
	//db.Exec("truncate auth")

	// 初始化
	//db.Exec("insert into admin_user (id, user_name, password, real_name, tel, status) values (1, 'admin', '$2a$10$057uuCLoKja2J04GLuWl1eNnwQtS7HxvookpbBa0thTHq7/fIaNF6', 'joy', '13054174174', 1)")

	//db.Exec("insert into role (id, pid, name, auth) values (1, 0, '超级管理员', '10,11,1100,1101,110000,110001,110002,110003,110004,110005,110100,110101,110102,110103,110104,110105')")
	//db.Exec("insert into role (id, pid, name, auth) values (2, 1, '系统维护管理员', '10,11,1100,1101,110000,110001,110002,110003,110004,110005,110100,110101,110102,110103,110104,110105')")

	//db.Exec("insert into casbin_rule (p_type, v0, v1) values ('g', 'user:1', 'role:1')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role', 'get')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role/:id', 'get')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role', 'post')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role/:id', 'put')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role/:id', 'delete')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser', 'get')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser/:id', 'get')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser', 'post')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser/:id', 'put')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser/:id', 'delete')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/auth/role', 'get')")
	//db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/auth/tree', 'get')")

	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (10, 0, '首页', 1, '', '', 'index')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (11, 0, '后台管理', 1, '', '', 'admin')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (1100, 11, '角色', 1, '', '', 'admin-role')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (1101, 11, '用户', 1, '', '', 'admin-user')")
	//
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110000, 1100, '获取权限树', 0, '/auth/tree', 'get', 'admin-user-auth-tree')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110001, 1100, '角色列表', 0, '/role', 'get', 'admin-role-list')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110002, 1100, '角色详情', 0, '/role/:id', 'get', 'admin-role-detail')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110003, 1100, '角色添加', 0, '/role', 'post', 'admin-role-add')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110004, 1100, '角色修改', 0, '/role/:id', 'put', 'admin-role-edit')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110005, 1100, '角色删除', 0, '/role/:id', 'delete', 'admin-role-del')")
	//
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110100, 1101, '获取角色树', 0, '/auth/role', 'get', 'admin-user-role-tree')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110101, 1101, '用户列表', 0, '/adminUser', 'get', 'admin-user-list')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110102, 1101, '用户详情', 0, '/adminUser/:id', 'get', 'admin-user-detail')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110103, 1101, '用户添加', 0, '/adminUser', 'post', 'admin-user-add')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110104, 1101, '用户修改', 0, '/adminUser/:id', 'put', 'admin-user-edit')")
	//db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110105, 1101, '用户删除', 0, '/adminUser/:id', 'delete', 'admin-user-del')")
}
