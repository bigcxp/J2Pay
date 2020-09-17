package model

import (
	"j2pay-server/pkg/setting"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var Db *gorm.DB

type Base struct{}
type FieldTrans map[string]string

func Setup() {
	var err error
	Db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
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

	Db.SingularTable(true)
	Db.LogMode(true)
	Db.SetLogger(&GormLogger{})
	Db.DB().SetMaxIdleConns(setting.MysqlConf.MaxIdle)
	Db.DB().SetMaxOpenConns(setting.MysqlConf.MaxActive)
	

	
 	AutoMigrate()

	// 设置程序启动参数 -init | -init=true
	if setting.Init {
		InitSql()
	}
}


// 通用分页获取偏移量
func GetOffset(page, pageSize int) int {
	if page <= 1 {
		return 0
	}
	return (page - 1) * pageSize
}


// 设置条件
func MultiWhere(where ...interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(where[0], where[1:]...)
	}
}


// 设置条件
func MultiOr(where ...interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(where[0], where[1:]...)
	}
}


// 自动创建修改表
func AutoMigrate() {
	Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '后台用户'").AutoMigrate(&AdminUser{})
	Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '角色'").AutoMigrate(&Role{})
	Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '权限'").AutoMigrate(&Auth{})
	Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT 'casbin policy 配置'").AutoMigrate(&CasbinRule{})
	Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '系统公告'").AutoMigrate(&SystemMessage{})
}

func InitSql() {
	// 清空
	Db.Exec("truncate admin_user")
	//Db.Exec("truncate role")
	//Db.Exec("truncate casbin_rule")
	//Db.Exec("truncate auth")

	// 初始化
	Db.Exec("insert into admin_user (id, user_name, password, real_name, tel, status) values (1, 'admin', '$2a$10$057uuCLoKja2J04GLuWl1eNnwQtS7HxvookpbBa0thTHq7/fIaNF6', 'joy', '13054174174', 1)")

	//Db.Exec("insert into role (id, pid, name, auth) values (1, 0, '超级管理员', '10,11,1100,1101,110000,110001,110002,110003,110004,110005,110100,110101,110102,110103,110104,110105')")
	//Db.Exec("insert into role (id, pid, name, auth) values (2, 1, '系统维护管理员', '10,11,1100,1101,110000,110001,110002,110003,110004,110005,110100,110101,110102,110103,110104,110105')")

	//Db.Exec("insert into casbin_rule (p_type, v0, v1) values ('g', 'user:1', 'role:1')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role', 'get')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role/:id', 'get')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role', 'post')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role/:id', 'put')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/role/:id', 'delete')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser', 'get')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser/:id', 'get')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser', 'post')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser/:id', 'put')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/adminUser/:id', 'delete')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/auth/role', 'get')")
	//Db.Exec("insert into casbin_rule (p_type, v0, v1, v2) values ('p', 'role:1', '/auth/tree', 'get')")

	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (10, 0, '首页', 1, '', '', 'index')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (11, 0, '后台管理', 1, '', '', 'admin')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (1100, 11, '角色', 1, '', '', 'admin-role')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (1101, 11, '用户', 1, '', '', 'admin-user')")
	//
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110000, 1100, '获取权限树', 0, '/auth/tree', 'get', 'admin-user-auth-tree')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110001, 1100, '角色列表', 0, '/role', 'get', 'admin-role-list')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110002, 1100, '角色详情', 0, '/role/:id', 'get', 'admin-role-detail')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110003, 1100, '角色添加', 0, '/role', 'post', 'admin-role-add')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110004, 1100, '角色修改', 0, '/role/:id', 'put', 'admin-role-edit')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110005, 1100, '角色删除', 0, '/role/:id', 'delete', 'admin-role-del')")
	//
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110100, 1101, '获取角色树', 0, '/auth/role', 'get', 'admin-user-role-tree')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110101, 1101, '用户列表', 0, '/adminUser', 'get', 'admin-user-list')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110102, 1101, '用户详情', 0, '/adminUser/:id', 'get', 'admin-user-detail')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110103, 1101, '用户添加', 0, '/adminUser', 'post', 'admin-user-add')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110104, 1101, '用户修改', 0, '/adminUser/:id', 'put', 'admin-user-edit')")
	//Db.Exec("insert into auth (id, pid, name, is_menu, api, action, ext) values (110105, 1101, '用户删除', 0, '/adminUser/:id', 'delete', 'admin-user-del')")
}
