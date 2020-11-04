package model

import (
	"j2pay-server/model/response"
	"j2pay-server/validate"
)

//提领 代发
type TWithdraw struct {
	ID           int64
	WithdrawType int       `gorm:"default:0;comment:' 类型 1 提币 2 代发'";json:"related_type"`                      // 类型 1 提币 2 代发
	SystemID     string    `gorm:"default:'';comment:'系统编号'";json:"system_id"`                                 // 系统编号
	MerchantID   string    `gorm:"default:'';comment:'商户订单编号'";json:"merchant_id"`                             // 商户订单编号
	ToAddress    string    `gorm:"default:'';comment:'提币地址'";json:"to_address"`                                // 提币地址
	Symbol       string    `gorm:"default:'';comment:'币种'";json:"symbol"`                                      //币种
	BalanceReal  float64   `gorm:"default:0;comment:'金额'";json:"balance_real"`                                 // 提币金额
	TxHash       string    `gorm:"default:'';comment:'提币tx hash'";json:"tx_hash"`                              // 提币tx hash
	Fee          float64   `gorm:"default:0;comment:'手续费'";json:"fee"`                                         //手续费
	CreateTime   int64     `gorm:"default:0;comment:'创建时间'";json:"create_time"`                                // 创建时间
	HandleStatus int64     `gorm:"default:0;comment:'处理状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败 '";json:"handle_status"` // 状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败
	HandleMsg    string    `gorm:"default:'';comment:'处理消息'";json:"handle_msg"`                                // 处理消息
	HandleTime   int64     `gorm:"default:'0';comment:'处理时间'";json:"handle_time"`                              // 处理时间
	Remark       string    `gorm:"default:'';comment:'备注';";json:"remark"`                                     //备注
	UserId       int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`
	AdminUser    AdminUser `json:"admin_user";gorm:"foreignkey:UserId"` //指定关联外键
}

// 管理端获取所有提领订单列表
func (p *TWithdraw) GetAllPick(page, pageSize int, where ...interface{}) (response.PickUpPage, error) {
	all := response.PickUpPage{
		Total:       p.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.PickList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Getdb().Model(&p).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.PickUpPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].RealName = GetUserByWhere("id = ?", v.UserId).RealName
	}

	return all, err
}

// 管理端获取所有代发订单列表
func (p *TWithdraw) GetAllSend(page, pageSize int, where ...interface{}) (response.SendPage, error) {
	all := response.SendPage{
		Total:       p.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.SendList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Getdb().Model(&p).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.SendPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].RealName = GetUserByWhere("id = ?", v.UserId).RealName
		all.Data[index].DelMoney = all.Data[index].Amount + all.Data[index].Fee
	}

	return all, err
}

// 商户端获取所有提领代发订单列表
func (p *TWithdraw) GetAll(page, pageSize int, where ...interface{}) (response.MerchantPickSendPage, error) {
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
	err := Getdb().Table("pick").Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.MerchantPickSendPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].RealName = GetUserByWhere("id = ?", v.UserId).RealName
		all.Data[index].DelMoney = all.Data[index].Amount + all.Data[index].Fee
		all.Data[index].GasFee = GetGasFeeDetail().EthFee
	}
	return all, err
}

// 获取所有提领订单数量
func (p *TWithdraw) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Getdb().Model(&p).Count(&count)
		return
	}
	Getdb().Model(&p).Where(where[0], where[1:]...).Count(&count)
	return
}

// 管理端根据ID获取提领订单详情
func (p *TWithdraw) GetPickDetail(id ...int) (res response.PickList, err error) {
	err = Getdb().Table("pick").
		Where("id = ?", p.ID).
		First(&res).
		Error
	res.RealName = GetUserByWhere("id = ?", res.UserId).RealName
	return
}

// 管理端根据ID获取提领订单详情
func (p *TWithdraw) GetSendDetail(id ...int) (res response.SendList, err error) {
	err = Getdb().Table("pick").
		Where("id = ?", p.ID).
		First(&res).
		Error
	res.RealName = GetUserByWhere("id = ?", res.UserId).RealName
	res.DelMoney = res.Amount + res.Fee
	return
}

// 商户端根据ID获取提领代发订单详情
func (p *TWithdraw) GetPickSendDetail(id ...int) (res response.MerchantPickList, err error) {

	err = Getdb().Table("pick").
		Where("id = ?", p.ID).
		First(&res).
		Error
	res.RealName = GetUserByWhere("id = ?", res.UserId).RealName
	res.DelMoney = res.Amount + res.Fee
	//gasFee 待完成·
	return
}

//取消代发 取消提领
func (p *TWithdraw) CancelPick(id, status int64) (err error) {
	tx := Getdb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	pick := GetPickByWhere("id = ?", id)
	err = tx.Model(&pick).
		Updates(TWithdraw{HandleStatus: status}).Error
	return
}

// 商户端创建提领或者代发订单
func (p *TWithdraw) Create() error {
	tx := Getdb().Begin()
	if err := tx.Create(p).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//获取提领总金额
func (p *TWithdraw) getAmount() float64 {
	var totalAmount float64
	all := response.MerchantPickSendPage{
		Data: []response.MerchantPickList{},
	}
	err := Getdb().Model(&p).Order("id desc").Where("status = ?", 2).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		totalAmount += validate.Decimal(v.Amount)
	}
	return validate.Decimal(totalAmount)
}

//总手续费
func (p *TWithdraw) getFee() float64 {
	var totalFee float64
	all := response.MerchantPickSendPage{
		Data: []response.MerchantPickList{},
	}
	err := Getdb().Model(&p).Order("id desc").Where("status = ?", 2).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		totalFee += validate.Decimal(v.Fee)
	}
	return validate.Decimal(totalFee)
}

// 根据条件获取详情
func GetPickByWhere(where ...interface{}) (pi TWithdraw) {
	Getdb().First(&pi, where...)
	return
}
