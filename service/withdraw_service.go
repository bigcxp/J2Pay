package service

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"j2pay-server/hcommon"
	"j2pay-server/heth"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
	"log"
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

// 提领 （提币 热钱包提币到用户自己账户）
func WithdrawAdd(with request.WithDrawAdd) (error, response.WithDrawRes) {
	defer casbin.ClearEnforcer()
	var c *gin.Context
	// 将币种小写
	with.Symbol = strings.ToLower(with.Symbol)
	// 根据登录信息获取组织id
	account, ok := c.Get("user")
	if !ok {
		return myerr.NewNormalValidateError("没有用户信息"), response.WithDrawRes{}
	}
	accountinfo := account.(*util.Claims)
	//获取组织信息
	user, err2 := model.GetUserByWhere("id = ?", accountinfo.UID)
	// 对比ip白名单
	if len(user.WhitelistIP) > 0 {
		if !strings.Contains(user.WhitelistIP, c.ClientIP()) {
			log.Println("no in ip list of: %s %s", user.RealName, c.ClientIP())
			return myerr.NewNormalValidateError("IP Limit"), response.WithDrawRes{}
		}
	}
	if err2 != nil {
		return err2, response.WithDrawRes{}
	}
	if user.ID == 0 {
		return myerr.NewNormalValidateError(" 没有该组织信息"), response.WithDrawRes{}
	}
	// 获取eth 信息
	tokenDecimalsMap := make(map[string]int64)
	ethSymbols := []string{heth.CoinSymbol}
	tokenDecimalsMap[heth.CoinSymbol] = 18
	// 获取所有eth代币币种
	tokenRows, err := model.SQLSelectTAppConfigTokenColAll()
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return myerr.NewNormalValidateError(err.Error()), response.WithDrawRes{}
	}
	for _, tokenRow := range tokenRows {
		tokenRow.TokenSymbol = strings.ToLower(tokenRow.TokenSymbol)
		ethSymbols = append(ethSymbols, tokenRow.TokenSymbol)
		tokenDecimalsMap[tokenRow.TokenSymbol] = tokenRow.TokenDecimals
	}
	// 验证金额
	tokenDecimals, ok := tokenDecimalsMap[with.Symbol]
	if !ok {
		return myerr.NewNormalValidateError("金额错误"), response.WithDrawRes{}
	}
	balanceObj, err := decimal.NewFromString(with.Balance)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return myerr.NewNormalValidateError(err.Error()), response.WithDrawRes{}
	}
	if balanceObj.LessThanOrEqual(decimal.NewFromInt(0)) {
		return myerr.NewNormalValidateError("金额不能小于等于0"), response.WithDrawRes{}
	}
	if balanceObj.Exponent() < -int32(tokenDecimals) {
		return myerr.NewNormalValidateError("金额格式错误"), response.WithDrawRes{}
	}
	if heth.IsStringInSlice(ethSymbols, with.Symbol) {
		// 验证组织的收款地址
		user.Address = strings.ToLower(user.Address)
		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
		if !re.MatchString(user.Address) {
			return myerr.NewNormalValidateError("组织收款地址格式错误"), response.WithDrawRes{}
		}
	} else {
		return myerr.NewNormalValidateError("暂不支持该币种"), response.WithDrawRes{}
	}
	now := time.Now().Unix()
	err = model.SQLCreateTWithdraw(
		&model.TWithdraw{
			UserId:       user.ID,
			SystemID:     util.RandString(12),
			ToAddress:    user.Address,
			Symbol:       with.Symbol,
			BalanceReal:  with.Balance,
			TxHash:       "",
			CreateTime:   now,
			HandleStatus: 0,
			HandleMsg:    "",
			HandleTime:   now,
			WithdrawType: hcommon.WithDraw,
			Remark:       with.Remark,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return myerr.NewNormalValidateError(err.Error()), response.WithDrawRes{}
	}
	return nil, response.WithDrawRes{with.Balance, user.Address}
}

// 代发 （代发 热钱包到外部钱包地址）
func SendAdd(with request.SendAdd) (error, response.SendRes) {
	defer casbin.ClearEnforcer()
	var c *gin.Context
	// 将币种小写
	with.Symbol = strings.ToLower(with.Symbol)
	// 根据登录信息获取组织id
	account, ok := c.Get("user")
	if !ok {
		return myerr.NewNormalValidateError("没有用户信息"), response.SendRes{}
	}
	accountinfo := account.(*util.Claims)
	//获取组织信息
	user, err2 := model.GetUserByWhere("id = ?", accountinfo.UID)
	if err2 != nil {
		return err2, response.SendRes{}
	}
	// 对比ip白名单
	if len(user.WhitelistIP) > 0 {
		if !strings.Contains(user.WhitelistIP, c.ClientIP()) {
			log.Println("no in ip list of: %s %s", user.RealName, c.ClientIP())
			return myerr.NewNormalValidateError("IP Limit"), response.SendRes{}
		}
	}
	if user.ID == 0 {
		return myerr.NewNormalValidateError(" 没有该组织i信息"), response.SendRes{}
	}
	//查询订单号否重复
	if hasName, _ := model.GetPickByWhere("user_name = ?", user.UserName); hasName.ID > 0 {
		return myerr.NewDbValidateError("订单号已存在"), response.SendRes{}
	}
	// 获取eth 信息
	tokenDecimalsMap := make(map[string]int64)
	ethSymbols := []string{heth.CoinSymbol}
	tokenDecimalsMap[heth.CoinSymbol] = 18
	// 获取所有eth代币币种
	tokenRows, err := model.SQLSelectTAppConfigTokenColAll()
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return myerr.NewNormalValidateError(err.Error()), response.SendRes{}
	}
	for _, tokenRow := range tokenRows {
		tokenRow.TokenSymbol = strings.ToLower(tokenRow.TokenSymbol)
		ethSymbols = append(ethSymbols, tokenRow.TokenSymbol)
		tokenDecimalsMap[tokenRow.TokenSymbol] = tokenRow.TokenDecimals
	}
	// 验证金额
	tokenDecimals, ok := tokenDecimalsMap[with.Symbol]
	if !ok {
		return myerr.NewNormalValidateError("金额错误"), response.SendRes{}
	}
	balanceObj, err := decimal.NewFromString(with.Balance)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return myerr.NewNormalValidateError(err.Error()), response.SendRes{}
	}
	if balanceObj.LessThanOrEqual(decimal.NewFromInt(0)) {
		return myerr.NewNormalValidateError("金额不能小于等于0"), response.SendRes{}
	}
	if balanceObj.Exponent() < -int32(tokenDecimals) {
		return myerr.NewNormalValidateError("金额格式错误"), response.SendRes{}
	}
	if heth.IsStringInSlice(ethSymbols, with.Symbol) {
		// 验证组织的收款地址
		user.Address = strings.ToLower(user.Address)
		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
		if !re.MatchString(user.Address) {
			return myerr.NewNormalValidateError("组织收款地址格式错误"), response.SendRes{}
		}
	} else {
		return myerr.NewNormalValidateError("暂不支持该币种"), response.SendRes{}
	}
	now := time.Now().Unix()
	err = model.SQLCreateTWithdraw(
		&model.TWithdraw{
			UserId:       user.ID,
			MerchantID:   with.OrderCode,
			SystemID:     util.RandString(12),
			ToAddress:    user.Address,
			Symbol:       with.Symbol,
			BalanceReal:  with.Balance,
			TxHash:       "",
			CreateTime:   now,
			HandleStatus: 0,
			HandleMsg:    "",
			HandleTime:   now,
			WithdrawType: hcommon.Send,
			Remark:       with.Remark,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return myerr.NewNormalValidateError(err.Error()), response.SendRes{}
	}
	return nil, response.SendRes{with.OrderCode,with.Balance, user.Address}

}

//取消提领 代发
func CancelPick(send request.SendEdit) error {
	defer casbin.ClearEnforcer()
	p := model.TWithdraw{}
	//逻辑
	return p.CancelPick(int64(send.ID), send.Status)
}
