package model

import (
	"j2pay-server/model/response"
	"j2pay-server/pkg/logger"
	"strconv"
	"strings"
)

type AdminUser struct {
	Id       int
	UserName string `gorm:"unique;comment:'用户名'"`
	Tel      string `gorm:"unique;default:'';comment:'手机号'"`
	Password string `gorm:"comment:'密码'"`
	RealName string `gorm:"default:'';comment:'真实姓名';"`
	Status   int8   `gorm:"default:1;comment:'状态 1:正常 0:停封'"`
}

// 获取所有后台用户
func (u *AdminUser) GetAll(page, pageSize int, where ...interface{}) (response.AdminUserPage, error) {
	all := response.AdminUserPage{
		Total:       u.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.AdminUserList{},
	}
	offset := GetOffset(page, pageSize)
	err := Db.Table("admin_user").
		Limit(pageSize).
		Offset(offset).
		Find(&all.Data, where...).Error
	if err != nil {
		return response.AdminUserPage{}, err
	}
	return all, err
}

// 根据ID获取用户详情
func (u *AdminUser) Detail(id ...int) (res response.AdminUserList, err error) {
	searchId := u.Id
	if len(id) > 0 {
		searchId = id[0]
	}
	err = Db.Table("admin_user").
		Where("id = ?", searchId).
		First(&res).
		Error
	return
}

// 创建
func (u *AdminUser) Create(roles []int) error {
	tx := Db.Begin()
	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, v := range roles {
		err := tx.Create(&CasbinRule{
			PType: "g",
			V0:    "user:" + strconv.Itoa(u.Id),
			V1:    "role:" + strconv.Itoa(v),
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 编辑用户
func (u *AdminUser) Edit(roles []int) error {
	tx := Db.Begin()
	updateInfo := map[string]interface{}{
		"user_name": u.UserName,
		"real_name": u.RealName,
		"status":    u.Status,
		"tel":       u.Tel,
	}
	if u.Password != "" {
		updateInfo["password"] = u.Password
	}
	if err := tx.Model(&AdminUser{Id: u.Id}).
		Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(CasbinRule{}, "p_type = 'g' and v0 = ?", "user:"+strconv.Itoa(u.Id)).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range roles {
		err := tx.Create(&CasbinRule{
			PType: "g",
			V0:    "user:" + strconv.Itoa(u.Id),
			V1:    "role:" + strconv.Itoa(v),
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 删除用户
func (u *AdminUser) Del() error {
	tx := Db.Begin()
	if err := tx.Delete(u, "id = ?", u.Id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(CasbinRule{}, "p_type = 'g' and v0 = ?", "user:"+strconv.Itoa(u.Id)).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 根据条件获取用户详情
func GetUserByWhere(where ...interface{}) (au AdminUser) {
	Db.First(&au, where...)
	return
}

// 获取所有后台用户数量
func (u *AdminUser) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Db.Model(&u).Count(&count)
		return
	}
	Db.Model(&u).Where(where[0], where[1:]...).Count(&count)
	return
}

// 根据用户 Id 获取所属角色
func GetUserRole(userId int) (userRoles []response.CasRole) {
	roles := GetAllRole()
	mappings := GetUserRoleMapping()
	_, ok := mappings[userId]
	if !ok {
		return
	}
	for _, role := range mappings[userId] {
		if _, ok := roles[role]; !ok {
			logger.Logger.Error("角色获取错误: user_id = ", userId)
			continue
		}
		userRoles = append(userRoles, roles[role])
	}
	return
}

// 根据用户 Id 获取权限
func GetUserAuth(userId int) (auth []Auth) {
	var roleId []int
	role := GetUserRole(userId)
	for _, v := range role {
		roleId = append(roleId, v.Id)
	}
	var dbRole []Role
	var whereAuthId []string
	Db.Model(Role{}).Select("auth").Find(&dbRole, "id in (?)", roleId)
	for _, v := range dbRole {
		whereAuthId = append(whereAuthId, v.Auth)
	}
	Db.Find(&auth, "id in (?)", strings.Split(strings.Join(whereAuthId, ","), ","))
	return
}

func GetAllUser() (mapping map[int]response.UserNames) {
	var users []response.UserNames
	mapping = make(map[int]response.UserNames)
	Db.Table("admin_user").Select("id,user_name").Order("id desc").Find(&users)
	for _, user := range users {
		mapping[user.Id] = user
	}
	return
}

// 根据条件获取多个角色
func GetUsersByWhere(where ...interface{}) (res []AdminUser, err error) {
	err = Db.Find(&res, where...).Error
	return
}
