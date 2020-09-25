package model

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/validate"
	"time"
)

type Order struct {
	gorm.Model
	OrderCode      string    `gorm:"default:'';comment:'商户订单编号';"json:"order_code"`
	Amount         float64   `gorm:"default:0;comment:'金额';";json:"amount"`
	ReceiptAmount  float64   `gorm:"default:0;comment:'实收金额';";json:"receipt_amount"`
	Fee            float64   `gorm:"default:0;comment:'手续费';";json:"fee"`
	ReturnAmount   float64   `gorm:"default:0;comment:'退款金额';";json:"return_amount"`
	MerchantAmount float64   `gorm:"default:0;comment:'商户实收金额';";json:"merchant_amount"`
	FinishTime     time.Time `gorm:"comment:'完成时间';";json:"finishTime"`
	TXID           string    `gorm:"default:'';comment:'交易哈希';";json:"txid"`
	Remark         string    `gorm:"default:'';commit:'备注';";json:"remark"`
	ChargeAddress  string    `gorm:"default:'';commit:'收款地址';";json:"charge_address"`
	UserId         int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`
	AdminUser      AdminUser `json:"admin_user";gorm:"foreignkey:UserId"` //指定关联外键
	Status         int       `gorm:"default:1;comment:'状态 -1：收款中，1：已完成，2：异常，3：退款等待中，4：退款中，5：退款失败，6：已退款，7：：已过期';";json:"status"`
}

//获取所有订单列表
func (o *Order) GetAll(page, pageSize int, where ...interface{}) (OrderPage, error) {
	all := OrderPage{
		Total:          o.GetCount(where...),
		PerPage:        pageSize,
		CurrentPage:    page,
		TotalAmount:    o.getTotalAmount(),
		ReallyAmount:   o.getReceiptAmount(),
		MerchantAmount: o.getMerchantAmount(),
		TotalFee:       o.getTotalFee(),
		Data:           []Order{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Db.Model(&o).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return OrderPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].AdminUser = GetUserByWhere("id = ?", v.UserId)
	}
	return all, err
}

// 获取所有订单数量
func (o *Order) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Db.Model(&o).Count(&count)
		return
	}
	Db.Model(&o).Where(where[0], where[1:]...).Count(&count)
	return
}

// 根据ID获取订单详情
func (o *Order) GetDetail(id ...int) (res Order, err error) {
	searchId := o.ID
	if len(id) > 0 {
		searchId = uint(id[0])
	}
	err = Db.Table("oeder").
		Where("id = ?", searchId).
		First(&res).
		Error
	adminUser := GetUserByWhere("id = ?", res.UserId)
	res.AdminUser = adminUser
	return
}

// 创建订单
func (o *Order) Create() error {
	tx := Db.Begin()
	o.CreatedAt = time.Now()
	o.FinishTime = time.Now()
	if err := tx.Create(o).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 根据条件获取订单详情
func GetOrderByWhere(where ...interface{}) (o Order) {
	Db.First(&o, where...)
	return
}

//商户订单总金额
func (o *Order) getTotalAmount() float64{
	var totalAmount float64
	all := OrderPage{
		Data:        []Order{},
	}
	err := Db.Model(&o).Order("id desc").Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		totalAmount += validate.Decimal(v.Amount)
	}
	return validate.Decimal(totalAmount)
}



//总商户总实收金额
func (o *Order) getMerchantAmount() float64{
	var merchantAmount float64
	all := OrderPage{
		Data:        []Order{},
	}
	err := Db.Model(&o).Order("id desc").Where("status = ?", 1).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		merchantAmount += validate.Decimal(v.MerchantAmount)
	}
	return validate.Decimal(merchantAmount)
}

//总实收金额
func (o *Order) getReceiptAmount() float64{
	var receiptAmount float64
	all := OrderPage{
		Data:        []Order{},
	}
	err := Db.Model(&o).Order("id desc").Where("status = ?", 1).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		receiptAmount += validate.Decimal(v.ReceiptAmount)
	}
	return validate.Decimal(receiptAmount)
}

//总手续费
func (o *Order) getTotalFee() float64 {
	var totalFee float64
	all := OrderPage{
		Data:        []Order{},
	}
	err := Db.Model(&o).Order("id desc").Where("status = ?", 1).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		totalFee += validate.Decimal(v.Fee)
	}
	return validate.Decimal(totalFee)
}