package model

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/model/response"
	"j2pay-server/validate"
	"time"
)

//提领订单
type Pick struct {
	gorm.Model
	IdCode      string    `gorm:"default:'';comment:'系统编号';";json:"id_code"`
	OrderCode   string    `gorm:"default:'';comment:'商户订单编号';"json:"order_code"`
	Amount      float64   `gorm:"default:0;comment:'金额';";json:"amount"`
	FinishTime  time.Time `gorm:"comment:'完成时间';";json:"finishTime"`
	TXID        string    `gorm:"default:'';comment:'交易信息';";json:"txid"`
	Fee         float64   `gorm:"default:0;comment:'手续费';";json:"fee"`
	Type        int       `gorm:"default:1;comment:'类型 1：代发，0：收款';";json:"type"`
	Remark      string    `gorm:"default:'';commit:'提领备注';";json:"remark"`
	PickAddress string    `gorm:"default:'';commit:'提领地址';";json:"pick_address"`
	Status      int       `gorm:"default:1;comment:'状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败';";json:"status"`
	UserId      int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`
	AdminUser   AdminUser `json:"admin_user";gorm:"foreignkey:UserId"` //指定关联外键
}

// 管理端获取所有提领订单列表
func (p *Pick) GetAllPick(page, pageSize int, where ...interface{}) (response.PickUpPage, error) {
	all := response.PickUpPage{
		Total:       p.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.PickList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Db.Model(&p).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.PickUpPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].RealName = GetUserByWhere("id = ?", v.UserId).RealName
	}

	return all, err
}
// 管理端获取所有代发订单列表
func (p *Pick) GetAllSend(page, pageSize int, where ...interface{}) (response.SendPage, error) {
	all := response.SendPage{
		Total:       p.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.SendList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Db.Model(&p).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.SendPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].RealName = GetUserByWhere("id = ?", v.UserId).RealName
		all.Data[index].DelMoney = all.Data[index].Amount+all.Data[index].Fee
	}

	return all, err
}
// 商户端获取所有提领代发订单列表
func (p *Pick) GetAll(page, pageSize int, where ...interface{}) (response.MerchantPickSendPage, error) {
	all := response.MerchantPickSendPage{
		Total:       p.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		TotalFee:    p.getFee(),
		TotalReduce: validate.Decimal(p.getFee() + p.getAmount()),
		TotalAmount: p.getAmount(),
		Data:        []response.MerchantPickList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Db.Model(&p).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.MerchantPickSendPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].RealName = GetUserByWhere("id = ?", v.UserId).RealName
		all.Data[index].DelMoney = all.Data[index].Amount+all.Data[index].Fee
		//gasFee 待完成
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

// 管理端根据ID获取提领订单详情
func (p *Pick) GetPickDetail(id ...int) (res response.PickList, err error) {
	searchId := p.ID
	if len(id) > 0 {
		searchId = uint(id[0])
	}
	err = Db.Table("pick").
		Where("id = ?", searchId).
		First(&res).
		Error
	res.RealName = GetUserByWhere("id = ?", res.UserId).RealName
	return
}

// 管理端根据ID获取提领订单详情
func (p *Pick) GetSendDetail(id ...int) (res response.SendList, err error) {
	searchId := p.ID
	if len(id) > 0 {
		searchId = uint(id[0])
	}
	err = Db.Table("pick").
		Where("id = ?", searchId).
		First(&res).
		Error
	res.RealName = GetUserByWhere("id = ?", res.UserId).RealName
	res.DelMoney = res.Amount+res.Fee
	return
}
// 商户端根据ID获取提领代发订单详情
func (p *Pick) GetPickSendDetail(id ...int) (res response.MerchantPickList, err error) {
	searchId := p.ID
	if len(id) > 0 {
		searchId = uint(id[0])
	}
	err = Db.Table("pick").
		Where("id = ?", searchId).
		First(&res).
		Error
	res.RealName = GetUserByWhere("id = ?", res.UserId).RealName
	res.DelMoney = res.Amount+res.Fee
	//gasFee 待完成·
	return
}

//取消代发 取消提领
func (p *Pick)CancelPick(id,status int)  error{
	tx := Db.Begin()
	pick := GetPickByWhere("id = ?", id)
	if err := tx.Model(&pick).
		Updates(Pick{Status: status}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// 商户端创建提领或者代发订单
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

//获取提领总金额
func (p *Pick) getAmount() float64 {
	var totalAmount float64
	all := response.MerchantPickSendPage{
		Data: []response.MerchantPickList{},
	}
	err := Db.Model(&p).Order("id desc").Where("status = ?", 2).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		totalAmount += validate.Decimal(v.Amount)
	}
	return validate.Decimal(totalAmount)
}

//总手续费
func (p *Pick) getFee() float64 {
	var totalFee float64
	all := response.MerchantPickSendPage{
		Data: []response.MerchantPickList{},
	}
	err := Db.Model(&p).Order("id desc").Where("status = ?", 2).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		totalFee += validate.Decimal(v.Fee)
	}
	return validate.Decimal(totalFee)
}


// 根据条件获取详情
func GetPickByWhere(where ...interface{}) (pi Pick) {
	Db.First(&pi, where...)
	return
}

