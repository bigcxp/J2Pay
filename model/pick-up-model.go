package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Pick struct {
	gorm.Model
	IdCode      string    `gorm:"default:'';comment:'系统编号';";json:"id_code"`
	OrderCode   string    `gorm:"default:'';comment:'商户订单编号';"json:"order_code"`
	Amount      float64   `gorm:"default:0;comment:'金额';";json:"amount"`
	FinishTime  time.Time `gorm:"comment:'完成时间';";json:"finishTime"`
	TXID        string    `gorm:"default:'';comment:'交易信息';";json:"txid"`
	Fee         float64   `gorm:"default:0;comment:'手续费';";json:"fee"`
	Type        int       `gorm:"default:1;comment:'类型 1：代发，0：收款';";json:"type"`
	UserId      int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`
	Remark      string    `gorm:"default:'';commit:'提领备注';";json:"remark"`
	PickAddress string    `gorm:"default:'';commit:'提领地址';";json:"pick_address"`
	Status      int       `gorm:"default:1;comment:'状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败';";json:"status"`
	AdminUser   AdminUser `json:"admin_user";gorm:"foreignkey:UserId"` //指定关联外键
}

// 获取所有提领订单列表
func (p *Pick) GetAll(page, pageSize int, where ...interface{}) (PickUpPage, error) {
	all := PickUpPage{
		Total:       p.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []Pick{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	//err := Db.Debug().Preload("Pick").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	err := Db.Model(&p).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return PickUpPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].AdminUser = GetUserByWhere("id = ?", v.UserId)
	}

	return all, err
}

// 获取所有提领订单数量
func (p *Pick) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Db.Model(&p).Count(&count)
		return
	}
	Db.Model(&p).Where(where[0], where[1:]...).Count(&count)
	return
}

// 根据ID获取提领订单详情
func (p *Pick) GetDetail(id ...int) (res Pick, err error) {
	searchId := p.ID
	if len(id) > 0 {
		searchId = uint(id[0])
	}
	err = Db.Table("pick").
		Where("id = ?", searchId).
		First(&res).
		Error
	adminUser := GetUserByWhere("id = ?", res.UserId)
	res.AdminUser = adminUser
	return
}

// 创建提领代发订单
func (p *Pick) Create() error {
	tx := Db.Begin()
	p.CreatedAt = time.Now()
	p.FinishTime = time.Now()
	if err := tx.Create(p).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
