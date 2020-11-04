package service

import (
	"github.com/shopspring/decimal"
	"j2pay-server/hcommon"
	"j2pay-server/heth"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
	"regexp"
	"strings"
	"time"
)

// 商户端  提领订单列表 + 代发订单列表
func MerchantPickList(fromDate, toDate, code string, types, userId, status, page, pageSize int) (res response.MerchantPickSendPage, err error) {
	pick := model.TWithdraw{}
	if userId == 0 {
		return res, myerr.NewDbValidateError("没有此用户信息")
	}
	if status == 0 && fromDate == "" && toDate == "" && code == "" && types == 0 {
		res, err = pick.GetAll(page, pageSize, "user_id = ?", userId)
	} else {
		//将时间进行转换
		res, err = pick.GetAll(page, pageSize, "user_id = ? and status = ? or type = ? or real_name like ? or is_code like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", userId, status, types, code, fromDate, toDate)
	}

	return
}

// 管理端 提领订单列表
func PickUpList(fromDate, toDate string, status, types, userId, page, pageSize int) (res response.PickUpPage, err error) {
	pick := model.TWithdraw{}

	if status == 0 && fromDate == "" && toDate == "" && userId == 0 {
		res, err = pick.GetAllPick(page, pageSize, "type = ?", types)
	} else {
		//将时间进行转换
		res, err = pick.GetAllPick(page, pageSize, "type = ? and status = ?  or user_id = ?  or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", types, status, userId, fromDate, toDate)
	}
	return
}

// 管理端  代发订单列表
func SendList(fromDate, toDate, orderCode string, status, types, userId, page, pageSize int) (res response.SendPage, err error) {
	pick := model.TWithdraw{}
	if orderCode == "" && status == 0 && fromDate == "" && toDate == "" && userId == 0 {
		res, err = pick.GetAllSend(page, pageSize, "type = ?", types)
	} else {
		//将时间进行转换
		res, err = pick.GetAllSend(page, pageSize, "type = ? and status = ? or orderCode like ? or user_id = ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", types, status, orderCode, userId, fromDate, toDate)
	}
	return
}

// 商户端提领订单或者代发订单详情
func MerchantPickDetail(id int64) (res response.MerchantPickList, err error) {
	pick := model.TWithdraw{
		ID: id,
	}
	res, err = pick.GetPickSendDetail()
	if err != nil {
		return
	}
	return res, err
}

// 管理端提领订单详情
func PickDetail(id int64) (res response.PickList, err error) {
	pick := model.TWithdraw{
		ID: id,
	}
	res, err = pick.GetPickDetail()
	if err != nil {
		return
	}
	return res, err
}

// 管理端代发订单详情
func SendDetail(id int64) (res response.SendList, err error) {
	pick := model.TWithdraw{
		ID: id,
	}
	res, err = pick.GetSendDetail()
	if err != nil {
		return
	}
	return res, err
}

// 提领 （提币 热钱包到用户自己账户）
func PickAdd(pick request.PickAdd) (error, response.PickAddr) {
	defer casbin.ClearEnforcer()
	//如果是RMB或TWD 换算成USDT
	var amount float64
	switch pick.Currency {
	case "RMB":
		detail, err := TypeDetail(pick.Currency)
		if err != nil {
			return err, response.PickAddr{}
		}
		amount = detail.OriginalRate / pick.Amount
	case "TWB":
		detail, err := TypeDetail(pick.Currency)
		if err != nil {
			return err, response.PickAddr{}
		}
		amount = detail.OriginalRate / pick.Amount

	default:
		amount = pick.Amount
	}
	//验证金额和地址是否正确
	tokenDecimalsMap := make(map[string]int64)
	ethSymbols := []string{heth.CoinSymbol}
	tokenDecimalsMap[heth.CoinSymbol] = 18
	token := model.TAppConfigToken{}
	tokenRows, err := token.SQLSelectTAppConfigTokenColAll()
	if err != nil {
		return err, response.PickAddr{}
	}
	for _, tokenRow := range tokenRows {
		tokenRow.TokenSymbol = strings.ToLower(tokenRow.TokenSymbol)
		ethSymbols = append(ethSymbols, tokenRow.TokenSymbol)
		tokenDecimalsMap[tokenRow.TokenSymbol] = tokenRow.TokenDecimals
	}
	// 验证金额
	tokenDecimals, ok := tokenDecimalsMap["eth"]
	if !ok {
		return err, response.PickAddr{}
	}
	balanceObj := decimal.NewFromFloat(pick.Amount)
	if balanceObj.LessThanOrEqual(decimal.NewFromInt(0)) {
		return myerr.NewDbValidateError("提领金额不能为负数"), response.PickAddr{}
	}
	if balanceObj.Exponent() < -int32(tokenDecimals) {
		return myerr.NewDbValidateError("小数位有误"), response.PickAddr{}
	}
	if hcommon.IsStringInSlice(ethSymbols, "eth") {
		// 验证地址
		pick.PickAddress = strings.ToLower(pick.PickAddress)
		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
		if !re.MatchString(pick.PickAddress) {
			return myerr.NewDbValidateError("地址不符合"), response.PickAddr{}
		}
	} else {
		return myerr.NewDbValidateError("地址不支持"), response.PickAddr{}
	}
	//逻辑处理 待完善 ==》判断金额 如果足够 进行交易，返回交易结果 账户减余额
	//代发功能是否开启
	if user := model.GetUserByWhere("id = ?", pick.UserId); user.IsDai == 0 {
		return myerr.NewDbValidateError("未开启代发功能"), response.PickAddr{}
	}
	// 1.判断转账金额是否足够
	if user := model.GetUserByWhere("id = ?", pick.UserId); user.Balance < pick.Amount {
		return myerr.NewDbValidateError("余额不足"), response.PickAddr{}
	}
	//提领数量不能为负数
	if pick.Amount < 0 {
		return myerr.NewDbValidateError("提领数量不能为负数"), response.PickAddr{}
	}
	//根据类型计算手续费
	user := model.GetUserByWhere("id = ?", pick.UserId)
	var fee float64
	if user.OrderCharge != 0 && user.OrderType == 1 {
		fee = pick.Amount * user.OrderCharge
	} else if user.OrderCharge != 0 && user.OrderType == 0 {
		fee = user.OrderCharge
	} else {
		fee = user.OrderCharge
	}
	var status int64
	var handleMsg string
	//是否需要审核
	if user.Examine == 1 {
		status = hcommon.PickStatusWait
		handleMsg = "等待中"
	} else {
		status = hcommon.PickStatusDo
		handleMsg = "执行中"
	}
	p := model.TWithdraw{
		SystemID:     util.RandString(20),
		BalanceReal:  pick.Amount,
		TxHash:       "",
		Fee:          fee,
		UserId:       pick.UserId,
		Remark:       pick.Remark,
		HandleStatus: status,
		HandleMsg:    handleMsg,
		HandleTime:   0,
		CreateTime:   time.Now().Unix(),
		ToAddress:    pick.PickAddress,
		Symbol:       "eth",
		WithdrawType: pick.Type,
	}
	pickAddr := response.PickAddr{
		Amount:         pick.Amount,
		Address:        pick.PickAddress,
		Currency:       pick.Currency,
		CurrencyAmount: amount,
	}
	return p.Create(), pickAddr
}

// 代发 （提币 热钱包到外部钱包地址）
func SendAdd(send request.SendAdd) (error, response.PickAddr) {
	defer casbin.ClearEnforcer()
	//如果是RMB或TWD 换算成USDT
	var amount float64
	switch send.Currency {
	case "RMB":
		detail, err := TypeDetail(send.Currency)
		if err != nil {
			return err, response.PickAddr{}
		}
		amount = detail.OriginalRate * send.Amount
	case "TWB":
		detail, err := TypeDetail(send.Currency)
		if err != nil {
			return err, response.PickAddr{}
		}
		amount = detail.OriginalRate * send.Amount

	default:
		amount = send.Amount
	}
	//验证金额和地址是否正确
	tokenDecimalsMap := make(map[string]int64)
	ethSymbols := []string{heth.CoinSymbol}
	tokenDecimalsMap[heth.CoinSymbol] = 18
	token := model.TAppConfigToken{}
	tokenRows, err := token.SQLSelectTAppConfigTokenColAll()
	if err != nil {
		return err, response.PickAddr{}
	}
	for _, tokenRow := range tokenRows {
		tokenRow.TokenSymbol = strings.ToLower(tokenRow.TokenSymbol)
		ethSymbols = append(ethSymbols, tokenRow.TokenSymbol)
		tokenDecimalsMap[tokenRow.TokenSymbol] = tokenRow.TokenDecimals
	}
	// 验证金额
	tokenDecimals, ok := tokenDecimalsMap["eth"]
	if !ok {
		return err, response.PickAddr{}
	}
	balanceObj := decimal.NewFromFloat(send.Amount)
	if balanceObj.LessThanOrEqual(decimal.NewFromInt(0)) {
		return myerr.NewDbValidateError("提领金额不能为负数"), response.PickAddr{}
	}
	if balanceObj.Exponent() < -int32(tokenDecimals) {
		return myerr.NewDbValidateError("小数位有误"), response.PickAddr{}
	}
	if hcommon.IsStringInSlice(ethSymbols, "eth") {
		// 验证地址
		send.PickAddress = strings.ToLower(send.PickAddress)
		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
		if !re.MatchString(send.PickAddress) {
			return myerr.NewDbValidateError("地址不符合"), response.PickAddr{}
		}
	} else {
		return myerr.NewDbValidateError("地址不支持"), response.PickAddr{}
	}
	//订单编号不能重复
	if hasCode := model.GetPickByWhere("order_code = ?", send.OrderCode); hasCode.ID > 0 {
		return myerr.NewDbValidateError("商户订单编号重复"), response.PickAddr{}
	}
	// 判断转账金额是否足够
	if user := model.GetUserByWhere("id = ?", send.UserId); user.Balance < send.Amount {
		return myerr.NewDbValidateError("余额不足"), response.PickAddr{}
	}
	//根据类型计算手续费
	user := model.GetUserByWhere("id = ?", send.UserId)
	var fee float64
	if user.OrderCharge != 0 && user.OrderType == 1 {
		fee = send.Amount * user.OrderCharge
	} else if user.OrderCharge != 0 && user.OrderType == 0 {
		fee = user.OrderCharge
	} else {
		fee = user.OrderCharge
	}
	//是否需要审核
	var status int64
	var handleMsg string
	//是否需要审核
	if user.Examine == 1 {
		status = hcommon.PickStatusWait
		handleMsg = "等待中"
	} else {
		status = hcommon.PickStatusDo
		handleMsg = "执行中"
	}
	p := model.TWithdraw{
		SystemID:     util.RandString(20),
		MerchantID:   send.OrderCode,
		BalanceReal:  send.Amount,
		TxHash:       "",
		Fee:          fee,
		UserId:       send.UserId,
		Remark:       send.Remark,
		HandleStatus: status,
		HandleMsg:    handleMsg,
		HandleTime:   0,
		CreateTime:   time.Now().Unix(),
		ToAddress:    send.PickAddress,
		Symbol:       "eth",
		WithdrawType: send.Type,
	}
	sendAddr := response.PickAddr{
		OrderCode:      send.OrderCode,
		Amount:         send.Amount,
		Address:        send.PickAddress,
		Currency:       send.Currency,
		CurrencyAmount: amount,
	}
	return p.Create(), sendAddr

}

//取消提领 代发
func CancelPick(send request.SendEdit) error {
	defer casbin.ClearEnforcer()
	p := model.TWithdraw{}
	//逻辑
	return p.CancelPick(int64(send.ID), send.Status)
}
