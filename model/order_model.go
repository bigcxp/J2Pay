package model

import (
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/validate"
	"strconv"
)

//订单表
type Order struct {
	ID            int64
	IdCode        string  `gorm:"default:'';comment:'系统编号';"json:"id_code"`
	OrderCode     string  `gorm:"default:'';comment:'商户订单编号';"json:"order_code"`
	ShouldAmount  float64 `gorm:"default:0;comment:'应收金额';";json:"should_amount"`
	ReceiptAmount float64 `gorm:"default:0;comment:'实收金额';";json:"receipt_amount"`
	Fee           float64 `gorm:"default:0;comment:'手续费';";json:"fee"`
	ReturnAmount  float64 `gorm:"default:0;comment:'退款金额';";json:"return_amount"`
	CreateTime    int64   `gorm:"comment:'创建时间';";json:"create_time"`
	FinishTime    int64   `gorm:"comment:'完成时间';";json:"finishTime"`
	ExprireTime   int64   `gorm:"comment:'过期时间';";json:"exprireTime"`
	Remark        string  `gorm:"default:'';comment:'备注';";json:"remark"`
	Address       string  `gorm:"default:'';comment:'收款地址';";json:"charge_address"`
	UserId        int64   `gorm:"default:0;comment:'组织ID'";json:"user_id"`
	TransactionId string  `gorm:"default:0;comment:'交易明细系统编号';"json:"transaction_id"`
	Status        int     `gorm:"default:-1;comment:'状态 -1：收款中，1：已完成，2：异常，3：退款等待中，4：退款中，5：退款失败，6：已退款，7：：已过期';";json:"status"`
}

//获取所有订单列表
func (o *Order) GetAllMerchantOrder(page, pageSize int, where ...interface{}) (response.OrderPage, error) {
	all := response.OrderPage{
		Total:          o.GetCount(where...),
		PerPage:        pageSize,
		CurrentPage:    page,
		TotalAmount:    o.getTotalAmount(),
		ReallyAmount:   o.getReceiptAmount(),
		MerchantAmount: o.getMerchantAmount(),
		TotalFee:       o.getTotalFee(),
		Data:           []response.RealOrderList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := DB.Model(&o).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.OrderPage{}, err
	}
	for index, v := range all.Data {
		user, _ := GetUserByWhere("id = ?", v.UserId)
		all.Data[index].RealName = user.RealName
		//根据订单明细系统编号获取
		txErc20, err := SQLSelectTTxErc20ByOrderId(v.OrderCode)
		if err != nil {
			return response.OrderPage{}, err
		}
		if txErc20.ID != 0 {
			all.Data[index].OrderDetail.ReceiptAmount = txErc20.BalanceReal
			all.Data[index].OrderDetail.TXID = txErc20.TxID
			all.Data[index].OrderDetail.DetailedRecordId = txErc20.SystemID
		}

	}
	return all, err
}

// 获取所有订单数量
func (o *Order) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		DB.Model(&o).Count(&count)
		return
	}
	DB.Model(&o).Where(where[0], where[1:]...).Count(&count)
	return
}

// 根据ID获取订单详情
func (o *Order) GetDetail(id ...int) (res response.RealOrderList, err error) {
	searchId := o.ID
	if len(id) > 0 {
		searchId = int64(id[0])
	}
	err = DB.Table("order").
		Where("id = ?", searchId).
		First(&res).
		Error
	user, _ := GetUserByWhere("id = ?", res.UserId)
	res.RealName = user.RealName
	//根据订单明细系统编号获取
	txErc20, err := SQLSelectTTxErc20ByOrderId(res.OrderCode)
	if err != nil {
		return response.RealOrderList{}, err
	}
	res.OrderDetail.ReceiptAmount = txErc20.BalanceReal
	res.OrderDetail.TXID = txErc20.TxID
	res.OrderDetail.DetailedRecordId = txErc20.SystemID
	return
}

// 创建订单
func (o *Order) Create() error {
	tx := DB.Begin()
	if err := tx.Create(o).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

//修改订单
func (o *Order) UpdateOrder(order request.OrderEdit) error {
	tx := DB.Begin()
	orders := GetOrderByWhere("id = ?", order.ID)
	if err := tx.Model(&orders).
		Updates(Order{Status: order.Status, Address: order.Address, ShouldAmount: order.ShouldAmount, ExprireTime: order.ExprireTime}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

//根据交易明细系统编号绑定交易明细记录
func BindErc20(Erc20Code string, orderCode string) (err error) {
	tx := GetDb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Model(&Order{}).Where("order_code = ?", orderCode).Update("transaction_id", Erc20Code).Error
	return
}

// 根据条件获取订单详情
func GetOrderByWhere(where ...interface{}) (o Order) {
	DB.First(&o, where...)
	return
}

//商户订单总金额
func (o *Order) getTotalAmount() float64 {
	var totalAmount float64
	all := response.OrderPage{
		Data: []response.RealOrderList{},
	}
	err := DB.Model(&o).Order("id desc").Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		totalAmount += validate.Decimal(v.Amount)
	}
	return validate.Decimal(totalAmount)
}

//总商户总实收金额
func (o *Order) getMerchantAmount() float64 {
	var merchantAmount float64
	all := response.OrderPage{
		Data: []response.RealOrderList{},
	}
	err := DB.Model(&o).Order("id desc").Where("status = ?", 1).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		merchantAmount += validate.Decimal(v.MerchantAmount)
	}
	return validate.Decimal(merchantAmount)
}

//总实收金额
func (o *Order) getReceiptAmount() float64 {
	var receiptAmount float64
	all := response.OrderPage{
		Data: []response.RealOrderList{},
	}
	err := DB.Model(&o).Order("id desc").Where("status = ?", 1).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		remount, _ := strconv.ParseFloat(v.OrderDetail.ReceiptAmount, 64)
		receiptAmount += remount
	}
	return validate.Decimal(receiptAmount)
}

//总手续费
func (o *Order) getTotalFee() float64 {
	var totalFee float64
	all := response.OrderPage{
		Data: []response.RealOrderList{},
	}
	err := DB.Model(&o).Order("id desc").Where("status = ?", 1).Find(&all.Data).Error
	if err != nil {
		return 0
	}
	for _, v := range all.Data {
		totalFee += validate.Decimal(v.Fee)
	}
	return validate.Decimal(totalFee)
}
