package model

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/model/response"
	"time"
)

//实收明细记录
type DetailedRecord struct {
	gorm.Model
	IdCode        string  `gorm:"default:'';comment:'系统订单编号';"json:"id_code"`
	Amount        float64 `gorm:"default:0;comment:'金额';";json:"amount"`
	TXID          string  `gorm:"default:'';comment:'交易哈希';";json:"txid"`
	Remark        string  `gorm:"default:'';comment:'备注';";json:"remark"`
	FromAddress   string  `gorm:"default:'';comment:'发款地址';";json:"from_address"`
	ChargeAddress string  `gorm:"default:'';comment:'收款地址';";json:"charge_address"`
	Status        int     `gorm:"default:1;comment:'状态 1：未绑定，2：已绑定';";json:"status"`
	OrderId       uint     `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"order_id"`
	Order         Order   `json:"order";gorm:"foreignkey:UserId"` //指定关联外键
	UserId      int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`
	AdminUser   AdminUser `json:"admin_user";gorm:"foreignkey:UserId"` //指定关联外键
}

//获取所有实收明细列表
func (d *DetailedRecord) GetAll(page, pageSize int, where ...interface{}) (response.DetailedRecordPage, error) {
	all := response.DetailedRecordPage{
		Total:       d.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.DetailedList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Db.Model(&d).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.DetailedRecordPage{}, err
	}
	for index, v := range all.Data {
		//如果没绑定 则没有组织和商户订单号
		if v.Status == 2 {
			all.Data[index].OrderCode = GetOrderByWhere("id = ?", v.OrderId).OrderCode
			all.Data[index].RealName = GetUserByWhere("id =?", v.UserId).RealName
		}

	}
	return all, err
}

//新增实收明细记录
func (d *DetailedRecord) DetailedAdd() error {
	tx := Db.Begin()
	d.CreatedAt = time.Now()
	if err := tx.Create(d).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 根据ID获取实收记录详情
func (d *DetailedRecord) GetDetail(id ...int) (res response.DetailedList, err error) {
	searchId := d.ID
	if len(id) > 0 {
		searchId = uint(id[0])
	}
	err = Db.Table("detail").
		Where("id = ?", searchId).
		First(&res).
		Error
	//如果该记录已经绑定
	if d.Status == 2 {
		order := GetOrderByWhere("id = ?", d.IdCode)
		res.OrderCode = order.OrderCode
		user := GetUserByWhere("id = ?", d.UserId)
		res.RealName = user.RealName
	}
	return
}

//解绑 绑定
func (d *DetailedRecord) Binding(id uint, orderCode string, isBind int) error {
	tx := Db.Begin()
	detailedRecord := GetDetailByWhere("id = ?", id)
	//解绑
	if isBind == 1 {
		if err := tx.Model(&detailedRecord).
			Updates(DetailedRecord{Status: 1, OrderId: 0,UserId: 0}).Error; err != nil {
			tx.Rollback()
			return err
		}
		//绑定
	} else {
		order := GetOrderByWhere("order_code = ?", orderCode)
		if err := tx.Model(&detailedRecord).
			Updates(DetailedRecord{Status: 2, OrderId: order.ID,UserId: order.UserId}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 获取所有实收明细记录数量
func (d *DetailedRecord) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Db.Model(&d).Count(&count)
		return
	}
	Db.Model(&d).Where(where[0], where[1:]...).Count(&count)
	return
}

// 根据条件获取详情
func GetDetailByWhere(where ...interface{}) (de DetailedRecord) {
	Db.First(&de, where...)
	return
}
