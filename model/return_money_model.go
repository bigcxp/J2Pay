package model

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/model/response"

	"time"
)

type Return struct {
	gorm.Model
	SystemCode string    `gorm:"default:'';comment:'系统订单编号';"json:"system_code"`
	OrderCode  string    `gorm:"default:'';comment:'商户订单编号';"json:"order_code"`
	Amount     float64   `gorm:"default:0;comment:'金额';";json:"amount"`
	FinishTime time.Time `gorm:"comment:'完成时间';";json:"finishTime"`
	UserId     int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`
	AdminUser  AdminUser `json:"admin_user";gorm:"foreignkey:UserId"` //指定关联外键
	Status     int       `gorm:"default:1;comment:'状态 1：退款等待中，2：退款中，3：退款失败，4：已退款';";json:"status"`
}

//获取所有订单列表
func (r *Return) GetAll(page, pageSize int, where ...interface{}) (response.ReturnPage, error) {
	all := response.ReturnPage{
		Total:       r.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.ReturnList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Getdb().Model(&r).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.ReturnPage{}, err
	}
	for index, v := range all.Data {
		user, _ := GetUserByWhere("id = ?", v.UserId)
		all.Data[index].RealName = user.RealName
	}
	return all, err
}

// 获取所有退款订单数量
func (r *Return) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Getdb().Model(&r).Count(&count)
		return
	}
	Getdb().Model(&r).Where(where[0], where[1:]...).Count(&count)
	return
}

// 根据ID获取退款订单详情
func (r *Return) GetDetail(id ...int) (res response.ReturnList, err error) {
	searchId := r.ID
	if len(id) > 0 {
		searchId = uint(id[0])
	}
	err = Getdb().Table("return").
		Where("id = ?", searchId).
		First(&res).
		Error
	user, _ := GetUserByWhere("id = ?", res.UserId)
	res.RealName = user.RealName
	return
}

// 创建退款订单
func (r *Return) Create() error {
	tx := Getdb().Begin()
	r.CreatedAt = time.Now()
	if err := tx.Create(r).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 根据条件获取订单详情
func GetReturnByWhere(where ...interface{}) (r Return) {
	Getdb().First(&r, where...)
	return
}
