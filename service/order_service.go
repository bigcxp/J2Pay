package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/moremorefun/mcommon"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
	"strings"
)

// 订单列表
func OrderList(fromDate string, toDate string, status int, chargeAddress string, txid string, orderCode string, userId int, page, pageSize int) (res response.OrderPage, err error) {
	order := model.Order{}
	if userId == 0 {
		if chargeAddress == "" && status == 0 && fromDate == "" && toDate == "" && orderCode == "" && txid == "" {
			res, err = order.GetAllMerchantOrder(page, pageSize)
		} else {
			//将时间进行转换
			res, err = order.GetAllMerchantOrder(page, pageSize, "status = ? or charge_address like ? or order_code like ?  or txid like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", status, chargeAddress, orderCode, txid, fromDate, toDate)
		}
	} else {
		if chargeAddress == "" && status == 0 && fromDate == "" && toDate == "" && orderCode == "" && txid == "" {
			res, err = order.GetAllMerchantOrder(page, pageSize, "user_id = ?", userId)
		} else {
			//将时间进行转换
			res, err = order.GetAllMerchantOrder(page, pageSize, "user_id = ? and status = ? or charge_address like ? or order_code like ?  or txid like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", userId, status, chargeAddress, orderCode, txid, fromDate, toDate)
		}
	}

	return
}

// 订单详情
func OrderDetail(id uint) (res response.RealOrderList, err error) {
	order := model.Order{}
	order.ID = int64(id)
	res, err = order.GetDetail()
	if err != nil {
		return
	}
	return res, err
}

// 新增订单 (获取充币地址)
func OrderAdd(c *gin.Context,order request.OrderAdd) (error, response.UserAddr) {
	defer casbin.ClearEnforcer()
	//获取当前登录用户
	account, ok := c.Get("user")
	if !ok {
		return  myerr.NewNormalValidateError("没有用户信息"),response.UserAddr{}
	}
	//判断订单编号是否重复
	 hasCode ,err:= model.GetOrderByWhere("order_code = ?", order.OrderCode)
	if err != nil {
		return err,response.UserAddr{}
	}
	if hasCode.ID > 0 {
		return myerr.NewDbValidateError("商户订单编号重复！"), response.UserAddr{}
	}
	//判断金额不能为负数
	if order.Amount < 0 {
		return myerr.NewNormalValidateError("值必须大于等于0"), response.UserAddr{}
	}
	accountinfo := account.(*util.Claims)
	a := model.Account{}
	count, _ := a.AccountDetail(accountinfo.ID)
	if count.RID != 1 {
		//获取组织
		user1 ,_:= model.GetUserByWhere("id = ?", count.UID)
		// 对比ip白名单
		if len(user1.WhitelistIP) > 0 {
			if !strings.Contains(user1.WhitelistIP, c.ClientIP()) {
				log.Warnf("no in ip list of: %s %s", user1.RealName, c.ClientIP())
				return myerr.NewNormalValidateError("IP Limit"), response.UserAddr{}
			}
		}

		// 接口验证签名
		var body []byte
		if cb, ok := c.Get(gin.BodyBytesKey); ok {
			if cbb, ok := cb.([]byte); ok {
				body = cbb
			}
		}
		if body == nil {
			body, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				return myerr.NewNormalValidateError("获取body为空"), response.UserAddr{}
			}
			c.Set(gin.BodyBytesKey, body)
		}
		oldObj := gin.H{}
		err = json.Unmarshal(body, &oldObj)
		if err != nil {
			log.Warnf("req body error")
			return myerr.NewNormalValidateError("req body erro"), response.UserAddr{}
		}
		checkObj := gin.H{}
		for k, v := range oldObj {
			if k != "sign" {
				checkObj[k] = v
			}
		}
		checkSign := mcommon.WechatGetSign(user1.UserName, checkObj)
		if checkSign == "" || checkSign != order.Sign {
			log.Warnf("sign error of: %s", user1.RealName)
			return myerr.NewNormalValidateError("sign error of"), response.UserAddr{}
		}

		//用户是否开启收款功能
		if user1.IsCollection != 1 {
			return myerr.NewNormalValidateError("未开启收款功能"), response.UserAddr{}
		}
		//是否开启手动建单
		if user1.IsCreation != 1 {
			return myerr.NewNormalValidateError("未开启手动建单功能"), response.UserAddr{}
		}
		var fee float64
		if user1.OrderCharge != 0 && user1.OrderType == 1 {
			fee = order.Amount * user1.OrderCharge
		} else if user1.OrderCharge != 0 && user1.OrderType == 0 {
			fee = user1.OrderCharge
		} else {
			fee = user1.OrderCharge
		}
		//检测当前用户的收钱地址是否满足 随机查询一个状态为已完成的地址
		address, err := model.GetAddress(int(user1.ID))
		if err != nil {
			return myerr.NewNormalValidateError("用户收款地址不足"), response.UserAddr{}
		}
		o := model.Order{
			OrderCode:   order.OrderCode,
			IdCode:      util.RandString(20),
			Address:     address.UserAddress,
			Fee:         fee,
			UserId:      user1.ID,
			CreateTime:  order.Uts,
			ExprireTime: order.Uts + user1.UserLessTime,
			Remark:      order.Remark,
		}
		userAddr := response.UserAddr{
			OrderCode:      order.OrderCode,
			Amount:         order.Amount,
			Address:        address.UserAddress,
			ExprireTime:    order.Uts + user1.UserLessTime,
		}
		return o.Create(), userAddr
		//管理员添加订单
	} else {
		//获取组织
		user1 ,_:= model.GetUserByWhere("id = ?", order.UID)
		// 对比ip白名单
		if len(user1.WhitelistIP) > 0 {
			if !strings.Contains(user1.WhitelistIP, c.ClientIP()) {
				log.Warnf("no in ip list of: %s %s", user1.RealName, c.ClientIP())
				return myerr.NewNormalValidateError("IP Limit"), response.UserAddr{}
			}
		}
		// 接口验证签名
		var body []byte
		if cb, ok := c.Get(gin.BodyBytesKey); ok {
			if cbb, ok := cb.([]byte); ok {
				body = cbb
			}
		}
		if body == nil {
			body, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				return myerr.NewNormalValidateError("获取body为空"), response.UserAddr{}
			}
			c.Set(gin.BodyBytesKey, body)
		}
		oldObj := gin.H{}
		err = json.Unmarshal(body, &oldObj)
		if err != nil {
			log.Warnf("req body error")
			return myerr.NewNormalValidateError("req body erro"), response.UserAddr{}
		}
		checkObj := gin.H{}
		for k, v := range oldObj {
			if k != "sign" {
				checkObj[k] = v
			}
		}
		checkSign := mcommon.WechatGetSign(user1.UserName, checkObj)
		if checkSign == "" || checkSign != order.Sign {
			log.Warnf("sign error of: %s", user1.RealName)
			return myerr.NewNormalValidateError("sign error of"), response.UserAddr{}
		}

		//组织是否开启收款功能
		if user1.IsCollection != 1 {
			return myerr.NewNormalValidateError("未开启收款功能"), response.UserAddr{}
		}
		//组织是否开启手动建单
		if user1.IsCreation != 1 {
			return myerr.NewNormalValidateError("未开启手动建单功能"), response.UserAddr{}
		}
		var fee float64
		if user1.OrderCharge != 0 && user1.OrderType == 1 {
			fee = order.Amount * user1.OrderCharge
		} else if user1.OrderCharge != 0 && user1.OrderType == 0 {
			fee = user1.OrderCharge
		} else {
			fee = user1.OrderCharge
		}
		//检测当前用户的收钱地址是否满足 随机查询一个状态为已完成的地址
		address, err := model.GetAddress(int(user1.ID))
		if err != nil {
			return myerr.NewNormalValidateError("用户收款地址不足"), response.UserAddr{}
		}
		o := model.Order{
			OrderCode:   order.OrderCode,
			IdCode:      util.RandString(20),
			Address:     address.UserAddress,
			Fee:         fee,
			UserId:      user1.ID,
			CreateTime:  order.Uts,
			ExprireTime: order.Uts + user1.UserLessTime,
			Remark:      order.Remark,
		}
		userAddr := response.UserAddr{
			OrderCode:      order.OrderCode,
			Amount:         order.Amount,
			Address:        address.UserAddress,
			ExprireTime:    order.Uts + user1.UserLessTime,
		}
		return o.Create(), userAddr
	}
}

//修改订单
func OrderEdit(order request.OrderEdit) error {
	defer casbin.ClearEnforcer()
	o := model.Order{}
	o.ID = int64(uint(order.ID))
	return o.UpdateOrder(order)
}

//订单通知
