package model

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/model/response"

	"time"
)

type Fee struct {
	gorm.Model
	Amount     float64   `gorm:"default:0;comment:'金额';";json:"amount"`
	FinishTime time.Time `gorm:"comment:'完成时间';";json:"finishTime"`
	UserId     int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`
	AdminUser  AdminUser `json:"admin_user";gorm:"foreignkey:UserId"` //指定关联外键
	Status     int       `gorm:"default:1;comment:'状态 1：等待中，2：已完成';";json:"status"`
}

//获取所有手续费列表
func (f *Fee) GetAll(page, pageSize int, where ...interface{}) (response.FeePage, error) {
	all := response.FeePage{
		Total:       f.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.FeeList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := DB.Model(&f).Order("id desc").Limit(pageSize).Offset(offset).Select("user.real_name, fee.*").Joins("inner join admin_user as user on fee.user_id = user.id").Find(&all.Data, where...).Error
	if err != nil {
		return response.FeePage{}, err
	}
	return all, err
}

//手续费结账
func (f *Fee) Settlement(id uint) error {
	tx := DB.Begin()
	fee := GetFeeByWhere("id = ?", id)
	if err := tx.Model(&fee).
		Updates(Fee{Status: 2}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

//创建手续费记录
func (f *Fee) Create() error {
	tx := DB.Begin()
	f.CreatedAt = time.Now()
	if err := tx.Create(f).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// 获取所有手续费列表数量
func (f *Fee) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		DB.Model(&f).Count(&count)
		return
	}
	DB.Model(&f).Where(where[0], where[1:]...).Count(&count)
	return

}

// 根据条件获取详情
func GetFeeByWhere(where ...interface{}) (fe Fee) {
	DB.First(&fe, where...)
	return
}
