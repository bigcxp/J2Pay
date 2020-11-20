package model

import (
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/validate"
	"strconv"
)

//提领 代发
type TWithdraw struct {
	ID           int64
	WithdrawType int       `gorm:"default:0;comment:' 类型 1 零钱整理 2 提币 3 代发'";json:"withdraw_type"`               //  类型 1 零钱整理 2 提币 3 代发
	SystemID     string    `gorm:"default:'';comment:'系统编号'";json:"system_id"`                                 // 系统编号
	MerchantID   string    `gorm:"default:'';comment:'商户订单编号'";json:"merchant_id"`                             // 商户订单编号
	ToAddress    string    `gorm:"default:'';comment:'提币地址'";json:"to_address"`                                // 提币地址
	Symbol       string    `gorm:"default:'';comment:'币种'";json:"symbol"`                                      //币种
	BalanceReal  string    `gorm:"default:0;comment:'金额'";json:"balance_real"`                                 // 提币金额
	TxHash       string    `gorm:"default:'';comment:'提币tx hash'";json:"tx_hash"`                              // 提币tx hash
	Fee          string    `gorm:"default:0;comment:'手续费'";json:"fee"`                                         //手续费
	CreateTime   int64     `gorm:"default:0;comment:'创建时间'";json:"create_time"`                                // 创建时间
	HandleStatus int64     `gorm:"default:0;comment:'处理状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败 '";json:"handle_status"` // 状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败
	HandleMsg    string    `gorm:"default:'';comment:'处理消息'";json:"handle_msg"`                                // 处理消息
	HandleTime   int64     `gorm:"default:'0';comment:'处理时间'";json:"handle_time"`                              // 处理时间
	Remark       string    `gorm:"default:'';comment:'备注';";json:"remark"`                                     //备注
	UserId       int64     `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`
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
	err := DB.Model(&p).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.PickUpPage{}, err
	}
	for index, v := range all.Data {
		user, _ := GetUserByWhere("id = ?", v.UserId)
		all.Data[index].RealName = user.RealName
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
	err := DB.Model(&p).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.SendPage{}, err
	}
	for index, v := range all.Data {
		user, _ := GetUserByWhere("id = ?", v.UserId)
		all.Data[index].RealName = user.RealName
		all.Data[index].DelMoney = all.Data[index].Amount + all.Data[index].Fee
	}

	return all, err
}

// 商户端获取所有提领代发订单列表
func (p *TWithdraw) GetAll(page, pageSize int, where ...interface{}) (response.MerchantPickSendPage, error) {
	fee, err2 := p.getFee()
	if err2 != nil {
		return response.MerchantPickSendPage{}, err2
	}
	amount, err2 := p.getAmount()
	if err2 != nil {
		return response.MerchantPickSendPage{}, err2
	}
	all := response.MerchantPickSendPage{
		Total:       p.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		TotalFee:    fee,
		TotalReduce: validate.Decimal(fee + amount),
		TotalAmount: amount,
		Data:        []response.MerchantPickList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := DB.Table("pick").Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.MerchantPickSendPage{}, err
	}
	for index, v := range all.Data {
		user, err := GetUserByWhere("id = ?", v.UserId)
		if err != nil {
			return response.MerchantPickSendPage{},err
		}
		all.Data[index].RealName = user.RealName
		all.Data[index].DelMoney = all.Data[index].Amount + all.Data[index].Fee
		detail, err := GetGasFeeDetail()
		if err != nil {
			return response.MerchantPickSendPage{}, err
		}
		all.Data[index].GasFee = detail.EthFee
	}
	return all, err
}

// 获取所有提领订单数量
func (p *TWithdraw) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		DB.Model(&p).Count(&count)
		return
	}
	DB.Model(&p).Where(where[0], where[1:]...).Count(&count)
	return
}

// 管理端根据ID获取提领订单详情
func (p *TWithdraw) GetPickDetail(id ...int) (res response.PickList, err error) {
	err = DB.Table("pick").
		Where("id = ?", p.ID).
		First(&res).
		Error
	user, _ := GetUserByWhere("id = ?", res.UserId)
	res.RealName = user.RealName
	return
}

// 管理端根据ID获取提领订单详情
func (p *TWithdraw) GetSendDetail(id ...int) (res response.SendList, err error) {
	err = DB.Table("pick").
		Where("id = ?", p.ID).
		First(&res).
		Error
	user, err := GetUserByWhere("id = ?", res.UserId)
	if err != nil {
		return
	}
	res.RealName = user.RealName
	res.DelMoney = res.Amount + res.Fee
	return
}

// 商户端根据ID获取提领代发订单详情
func (p *TWithdraw) GetPickSendDetail(id ...int) (res response.MerchantPickList, err error) {

	err = DB.Table("pick").
		Where("id = ?", p.ID).
		First(&res).
		Error
	user, err := GetUserByWhere("id = ?", res.UserId)
	if err != nil {
		return
	}
	res.RealName = user.RealName
	res.DelMoney = res.Amount + res.Fee
	//gasFee 待完成·
	return
}

//取消代发 取消提领
func (p *TWithdraw) CancelPick(id, status int64) (err error) {
	tx := DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()

		}
	}()
	pick,err := GetPickByWhere("id = ?", id)
	if err != nil {
		return err
	}
	err = tx.Model(&pick).
		Updates(TWithdraw{HandleStatus: status}).Error
	return
}

// 商户端创建提领或者代发订单
func SQLCreateTWithdraw(row *TWithdraw) error {
	if row == nil {
		return myerr.NewDbValidateError("no data")
	}
	tx := DB.Begin()
	if err := tx.Create(row).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//获取提领总金额
func (p *TWithdraw) getAmount() (float64,error) {
	var totalAmount int
	all := response.MerchantPickSendPage{
		Data: []response.MerchantPickList{},
	}
	err := DB.Model(&p).Order("id desc").Where("status = ?", 2).Find(&all.Data).Error
	if err != nil {
		return 0,err
	}
	for _, v := range all.Data {
		atoi, _ := strconv.Atoi(v.Amount)
		totalAmount +=atoi
	}
	return validate.Decimal(float64(totalAmount)),nil
}

//总手续费
func (p *TWithdraw) getFee() (float64,error) {
	var totalFee int
	all := response.MerchantPickSendPage{
		Data: []response.MerchantPickList{},
	}
	err := DB.Model(&p).Order("id desc").Where("status = ?", 2).Find(&all.Data).Error
	if err != nil {
		return 0,err
	}
	for _, v := range all.Data {
		fee, _ := strconv.Atoi(v.Fee)
		totalFee += fee
	}
	return validate.Decimal(float64(totalFee)),nil
}

// 根据条件获取详情
func GetPickByWhere(where ...interface{}) ( TWithdraw,error) {
	var pi TWithdraw
	err := DB.First(&pi, where...).Error
	return pi,err
}
