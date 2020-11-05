package model

import (
	"j2pay-server/model/response"
	"j2pay-server/pkg/logger"
	"strconv"
	"time"
)

type AdminUser struct {
	ID             int64
	UserName       string          `gorm:"unique;comment:'主账户名';"`
	RealName       string          `gorm:"unique;comment:'组织名称';"`
	WhitelistIP    string          `gorm:"default:'';comment:'ip白名单';"`
	CreateTime     int64           `gorm:"type:timestamp;comment:'创建时间';"`
	UpdateTime     int64           `gorm:"type:timestamp;comment:'更新时间';"`
	Balance        float64         `gorm:"default:0;comment:'组织余额';"`
	Address        string          `gorm:"default:'';comment:'商户地址';"`
	ReturnUrl      string          `gorm:"default:'';comment:'回传URL'"`
	DaiUrl         string          `gorm:"default:'';comment:'代发URL'"`
	Remark         string          `gorm:"default:'';comment:'备注'"`
	IsCollection   int             `gorm:"default:1;comment:'是否开启收款功能 1：是 0：否'"`
	IsCreation     int             `gorm:"default:1;comment:'是否开启手动建单 1：是 0：否'"`
	More           int             `gorm:"default:0;comment:'地址多单收款'"`
	OrderType      int             `gorm:"default:1;comment:'订单手续费类型 1：百分比 0：固定'"`
	OrderCharge    float64         `gorm:"default:0;comment:'订单手续费';"`
	ReturnType     int             `gorm:"default:1;comment:'退款手续费类型 1：百分比 0：固定'"`
	ReturnCharge   float64         `gorm:"default:0;comment:'退款手续费';"`
	IsDai          int             `gorm:"default:1;comment:'是否启用代发功能 1：是 0：否';"`
	PickType       int             `gorm:"default:1;comment:'提领手续费类型 1：百分比 0：固定'"`
	PickCharge     float64         `gorm:"default:0;comment:'提领手续费';"`
	DaiType        int             `gorm:"default:1;comment:'代发手续费类型 1：百分比 0：固定'"`
	DaiCharge      float64         `gorm:"default:0;comment:'代发手续费';"`
	IsGas          int             `gorm:"default:1;comment:'是否启用gas预估 1：是 0：否'"`
	Examine        float64         `gorm:"default:0;comment:'代发审核 1：是 0：否 ';"`
	DayTotalCount  float64         `gorm:"default:0;comment:'每日交易总量';"`
	MaxOrderCount  float64         `gorm:"default:0;comment:'最大交易数量';"`
	MinOrderCount  float64         `gorm:"default:0;comment:'最小交易数量';"`
	Limit          float64         `gorm:"default:0;comment:'结账限制';"`
	UserLessTime   int64           `gorm:"default:0;comment:'订单无效时间';"`
	TWithdraw      []TWithdraw     `gorm:"FOREIGNKEY:UserId;ASSOCIATION_FOREIGNKEY:Id"`
	SystemMessages []SystemMessage `gorm:"many2many:system_message_user;"`
}

// 获取所有后台组织
func (u *AdminUser) GetAll(page, pageSize int, where ...interface{}) (response.AdminUserPage, error) {
	all := response.AdminUserPage{
		Total:       u.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.AdminUserList{},
	}
	offset := GetOffset(page, pageSize)
	err := Getdb().Table("admin_user").
		Limit(pageSize).
		Offset(offset).
		Find(&all.Data, where...).Error
	if err != nil {
		return response.AdminUserPage{}, err
	}
	return all, err
}

// 获取组织所有系统公告
func (u *AdminUser) GetAllMessage(page, pageSize int, where ...interface{}) (response.AdminUserMessagePage, error) {
	user := &response.AdminUserMessageList{}
	Getdb().Where("id = ?", 1).Preload("system_message").Find(&user)
	logger.Logger.Println(user)

	all := response.AdminUserMessagePage{
		Total:       u.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.AdminUserMessageList{},
	}
	offset := GetOffset(page, pageSize)
	err := Getdb().Table("admin_user").
		Limit(pageSize).
		Offset(offset).
		Find(&all.Data, where...).Error
	if err != nil {
		return response.AdminUserMessagePage{}, err
	}
	return all, err
}

// 根据ID获取组织详情
func (u *AdminUser) Detail(id ...int64) (res response.AdminUserList, err error) {
	searchId := u.ID
	if len(id) > 0 {
		searchId = id[0]
	}
	err = Getdb().Table("admin_user").
		Where("id = ?", searchId).
		First(&res).
		Error
	return
}

// 创建组织
func (u *AdminUser) Create()(int64,error) {
	tx := Getdb().Begin()
	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		return 0,err
	}
	tx.Commit()
	return u.ID,nil
}

// 编辑组织
func (u *AdminUser) Edit() error {
	tx := Getdb().Begin()
	updateInfo := map[string]interface{}{
		"real_name":       u.RealName,
		"address":         u.Address,
		"dai_url":         u.DaiUrl,
		"return_url":      u.ReturnUrl,
		"whitelist_ip":    u.WhitelistIP,
		"is_collection":   u.IsCollection,
		"is_creation":     u.IsCreation,
		"more":            u.More,
		"order_type":      u.OrderType,
		"order_charge":    u.OrderCharge,
		"return_type":     u.ReturnType,
		"return_charge":   u.ReturnCharge,
		"is_dai":          u.IsDai,
		"dai_type":        u.DaiType,
		"dai_charge":      u.DaiCharge,
		"pick_type":       u.PickType,
		"pick_charge":     u.PickCharge,
		"is_gas":          u.IsGas,
		"examine":         u.Examine,
		"day_total_count": u.DayTotalCount,
		"max_order_count": u.MaxOrderCount,
		"min_order_count": u.MinOrderCount,
		"limit":           u.Limit,
		"user_less_time":  u.UserLessTime,
		"update_time":     time.Now().Unix(),
	}
	if err := tx.Model(&AdminUser{ID: u.ID}).
		Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 删除组织
func (u *AdminUser) Del() error {
	tx := Getdb().Begin()
	if err := tx.Delete(u, "id = ?", u.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(CasbinRule{}, "p_type = 'g' and v0 = ?", "user:"+strconv.Itoa(int(u.ID))).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 根据条件获取组织详情
func GetUserByWhere(where ...interface{}) (au AdminUser, err error) {
	err = Getdb().First(&au, where...).Error
	return
}

// 获取所有后台组织数量
func (u *AdminUser) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Getdb().Model(&u).Count(&count)
		return
	}
	Getdb().Model(&u).Where(where[0], where[1:]...).Count(&count)
	return
}

func GetAllUser() (mapping map[int]response.UserNames) {
	var users []response.UserNames
	mapping = make(map[int]response.UserNames)
	Getdb().Table("admin_user").Select("id,user_name").Order("id desc").Find(&users)
	for _, user := range users {
		mapping[user.Id] = user
	}
	return
}

// 根据条件获取多个角色
func GetUsersByWhere(where ...interface{}) (res []AdminUser, err error) {
	err = Getdb().Find(&res, where...).Error
	return
}
