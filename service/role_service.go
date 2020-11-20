package service

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"strconv"
	"strings"
)

// 角色列表
func RoleList(page, pageSize int) (response.RolePage, error) {
	role := model.Role{}
	return role.GetAll(page, pageSize)
}

// 角色详情
func RoleDetail(id int) (response.RoleList, error) {
	role := model.Role{ID: id}
	res, err := role.Detail()
	if err != nil {
		return res, err
	}
	auth := strings.Split(res.Auths, ",")
	for _, v := range auth {
		id, _ := strconv.Atoi(v)
		res.Auth = append(res.Auth, id)
	}
	baseAuth ,err:= model.GetAllBaseAuth("is_menu = 0 and id in (?)", res.Auth)
	if err != nil {
		return res,err
	}
	for _, v := range baseAuth {
		res.BaseAuth = append(res.BaseAuth, v.Id)
	}
	return res, err
}

// 添加角色
func RoleAdd(role request.RoleAdd) error {
	defer casbin.ClearEnforcer()
	r := model.Role{
		Name: role.Name,
	}
	// 判断是否有重复的角色名
	hasRole, err := model.GetRoleByWhere("name = ?", r.Name)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if hasRole.ID > 0 {
		return myerr.NewDbValidateError("角色名已存在")
	}

	// 严谨起见，从数据库获取权限
	all,err := model.GetAllBaseAuth("id in (?)", role.Auth)
	if err != nil {
		return err
	}
	allIds := make([]string, 0)
	for _, v := range all {
		allIds = append(allIds, strconv.Itoa(v.Id))
	}
	r.Auth = strings.Join(allIds, ",")

	// 创建角色
	return r.Create(all)
}

// 编辑角色
func RoleEdit(role request.RoleEdit) error {
	defer casbin.ClearEnforcer()
	r := model.Role{
		ID:   role.ID,
		Name: role.Name,
	}
	// 判断是否有重复的角色名
	hasRole, err := model.GetRoleByWhere("name = ? and id <> ?", r.Name, r.ID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if hasRole.ID > 0 {
		return myerr.NewDbValidateError("角色名已存在")
	}

	// 严谨起见，从数据库获取权限
	all ,err:= model.GetAllBaseAuth("id in (?)", role.Auth)
	if err != nil {
		return err
	}
	allIds := make([]string, 0)
	for _, v := range all {
		allIds = append(allIds, strconv.Itoa(v.Id))
	}
	r.Auth = strings.Join(allIds, ",")

	// 编辑角色
	return r.Edit(all)
}

// 删除角色
func RoleDel(id int) error {
	defer casbin.ClearEnforcer()
	role := model.Role{ID: id}

	// 查看用户是否使用该角色
	key := "role:" + strconv.Itoa(id)
	res, err := model.GetCasbinByWhere("p_type = 'g' and (v0 = ? or v1 = ?)", key, key)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if res.PType != "" {
		return myerr.NewDbValidateError("该角色已被使用，无法删除")
	}

	// 同步更新 casbin
	return role.Del()
}

// 获取角色树
func RoleTree(self int) (res []response.Roles, err error) {
	res, err = model.GetRoleTreeByWhere("id <> ?", self)
	if err != nil {
		return
	}
	return
}
