package model

import (
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/validate"
	"strconv"
	"strings"
	"time"
)

//账户实体
type Account struct {
	ID            int64
	UID           int64  `gorm:"default:0;comment:'组织ID'"`
	RID           int    `gorm:"default:0;comment:'角色ID'"`
	UserName      string `gorm:"unique;comment:'用户名';"`
	Password      string `gorm:"comment:'密码';"`
	Token         string `gorm:"default:'';comment:'Token'"`
	Secret        string `gorm:"comment:'google私钥';"`
	IsOpen        int    `gorm:"default:0;comment:'是否开启google双重验证 默认0：不开启 1：开启';"`
	QrcodeUrl     string `gorm:"comment:'google二维码图片地址';"`
	LastLoginTime int64  `gorm:"type:timestamp;comment:'最后登录时间';"`
	CreateTime    int64  `gorm:"default:0;comment:'创建时间'";json:"create_time"`
	UpdateTime    int64  `gorm:"default:0;comment:'修改时间'";json:"update_time"`
	Status        int    `gorm:"default:1;comment:'是否启用 1:正常 0:停封'"`
	IsMain        int    `gorm:"default:0;comment:'是否是主账号 1:是 0:否'"`
}

// 根据条件获取账户详情
func GetAccountByWhere(where ...interface{}) (ac Account,err error) {
	err = Getdb().First(&ac, where...).Error
	return
}

// 获取所有账户
func (a *Account) GetAll(page, pageSize int, where ...interface{}) (response.AccountPage, error) {
	all := response.AccountPage{
		Total:       a.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.AccountList{},
	}
	offset := GetOffset(page, pageSize)
	err := Getdb().Model(&Account{}).
		Limit(pageSize).
		Offset(offset).
		Find(&all.Data, where...).Error
	if err != nil {
		return response.AccountPage{}, err
	}
	return all, err
}

// 根据ID获取账户详情
func (a *Account) AccountDetail(id ...int64) (res response.AccountList, err error) {
	searchId := a.ID
	if len(id) > 0 {
		searchId = int64(id[0])
	}
	err = Getdb().Model(&Account{}).
		Where("id = ?", searchId).
		First(&res).
		Error
	return
}

// 创建账户
func (a *Account) Create(role int) error {
	tx := Getdb().Begin()
	if err := tx.Create(a).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&CasbinRule{
		PType: "g",
		V0:    "user:" + strconv.Itoa(int(a.ID)),
		V1:    "role:" + strconv.Itoa(role),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 编辑账户
func (a *Account) Edit(role int) error {
	tx := Getdb().Begin()
	updateInfo := map[string]interface{}{
		"status":      a.Status,
		"update_time": a.UpdateTime,
	}
	if err := tx.Model(&Account{ID: a.ID}).
		Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(CasbinRule{}, "p_type = 'g' and v0 = ?", "user:"+strconv.Itoa(int(a.ID))).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&CasbinRule{
		PType: "g",
		V0:    "user:" + strconv.Itoa(int(a.ID)),
		V1:    "role:" + strconv.Itoa(role),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//修改密码
func (a *Account) UpdatePassword(id int64, password string) (err error) {
	tx := Getdb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	user,_ := GetUserByWhere("id = ?", id)
	err = tx.Model(&user).
		Updates(Account{Password: password}).Error
	return
}

//是否开启google验证
func (a *Account) Google(google request.Google) (err error) {
	tx := Getdb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	user,_ := GetUserByWhere("id = ?", google.ID)
	isOpen := google.IsOpen
	if isOpen == 0 {
		err = tx.Model(&Account{ID: user.ID}).
			Updates(Account{IsOpen: 0}).Error
	} else {
		err = tx.Model(&AdminUser{ID: user.ID}).
			Updates(Account{IsOpen: 1}).Error
	}
	return
}

// 删除用户
func (a *Account) Del() error {
	tx := Getdb().Begin()
	if err := tx.Delete(a, "id = ?", a.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(CasbinRule{}, "p_type = 'g' and v0 = ?", "user:"+strconv.Itoa(int(a.ID))).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 根据用户 ID 获取所属角色
func GetAccountRole(userId int64) (userRoles response.CasRole) {
	roles := GetAllRole()
	mappings := GetAccountRoleMapping()
	_, ok := mappings[userId]
	if !ok {
		return
	}
	userRoles = roles[mappings[userId]]
	return
}

// 根据账户 ID 获取权限
func GetAccountAuth(Id int64) (auth []Auth) {
	role := GetAccountRole(Id)
	var dbRole Role
	var whereAuthId []string
	Getdb().Model(Role{}).Select("auth").First(&dbRole, "id = ?", role.ID)
	whereAuthId = append(whereAuthId, dbRole.Auth)
	Getdb().Find(&auth, "id in (?)", strings.Split(strings.Join(whereAuthId, ","), ","))
	return
}

func GetAllAccount() (mapping map[int]response.UserNames) {
	var users []response.UserNames
	mapping = make(map[int]response.UserNames)
	Getdb().Table("admin_user").Select("id,user_name").Order("id desc").Find(&users)
	for _, user := range users {
		mapping[user.Id] = user
	}
	return
}

// 编辑用户
func (u *Account) EditToken(username, token string) error {
	tx := Getdb().Begin()
	account ,_:= GetAccountByWhere("user_name = ?", username)
	if err := tx.Model(&account).
		Updates(Account{LastLoginTime: time.Now().Unix(), QrcodeUrl: validate.NewGoogleAuth().GetQrcodeUrl(username, account.Secret), Token: token}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 获取所有账户数量
func (a *Account) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Getdb().Model(&a).Count(&count)
		return
	}
	Getdb().Model(&a).Where(where[0], where[1:]...).Count(&count)
	return
}
