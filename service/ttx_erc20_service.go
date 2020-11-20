package service

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
	"time"
)

// 实收明细记录列表
func Erc20List(c *gin.Context,IdCode, address, txid, fromDate, toDate string, status, page, pageSize int) (res response.Erc20Page, err error) {
	ttxErc20 := model.TTxErc20{}
	//获取当前登录用户
	account, ok := c.Get("user")
	if !ok {
		return response.Erc20Page{}, myerr.NewNormalValidateError("没有用户信息")
	}
	accountinfo := account.(*util.Claims)
	a := model.Account{}
	count, _ := a.AccountDetail(accountinfo.ID)
	//如果是管理员
	if count.RID == 1 {
		if status == 0 && fromDate == "" && toDate == "" && txid == "" && address == "" && IdCode == "" {
			res, err = ttxErc20.GetAllErc20Detail(page, pageSize)
		} else {
			if fromDate == "" || toDate == ""{
				res, err = ttxErc20.GetAllErc20Detail(page, pageSize, "status = ? or system_id like ? or tx_id like ? or to_address like ?", status, IdCode, txid, address)
			} else {
				res, err = ttxErc20.GetAllErc20Detail(page, pageSize, "status = ? or system_id like ? or tx_id like ? or to_address like ? or create_time>=? or create_time <=?", status, IdCode, txid, address,util.TimeStr2Time(fromDate).Unix(), util.TimeStr2Time(toDate).Unix())
			}
		}
		if err != nil {
			return
		}
	} else {
		//获取组织信息
		user, err2 := model.GetUserByWhere("id = ?", accountinfo.UID)
		if err2 != nil {
			return response.Erc20Page{}, myerr.NewNormalValidateError("没有组织信息")
		}
		if status == 0 && fromDate == "" && toDate == "" && txid == "" && address == "" && IdCode == "" {
			res, err = ttxErc20.GetAllErc20Detail(page, pageSize, "user_id = ?", user.ID)
		} else {
			if fromDate == "" || toDate == ""{
				res, err = ttxErc20.GetAllErc20Detail(page, pageSize, "user_id  = ? and status = ? or system_id like ? or tx_id like ? or to_address like ?", status, IdCode, txid, address)
			} else {
				res, err = ttxErc20.GetAllErc20Detail(page, pageSize, "user_id  = ? and status = ? or system_id like ? or tx_id like ? or to_address like ? or create_time>=? or create_time <=?", status, IdCode, txid, address,util.TimeStr2Time(fromDate).Unix(), util.TimeStr2Time(toDate).Unix())
			}
		}
		if err != nil {
			return
		}
	}
	return
}

//新增实收明细记录
func Erc20Add(erc20add request.Erc20Add) error {
	defer casbin.ClearEnforcer()
	now := time.Now().Unix()
	token := model.TAppConfigToken{}
	selectBySymbol, err := token.SQLSelectBySymbol("usdt")
	if err != nil {
		return err
	}
	t := model.TTxErc20{
		TxID:         erc20add.TxID,
		TokenID:      selectBySymbol.ID,
		SystemID:     util.RandString(12),
		FromAddress:  erc20add.From,
		ToAddress:    erc20add.To,
		BalanceReal:  erc20add.Balance,
		Remark:       erc20add.Remark,
		CreateTime:   now,
		HandleStatus: 1,
		OrgStatus:    0,
	}
	// 1.判断txid是否重复
	 hasTx ,err:= model.GetErc20ByWhere("tx_id = ?", erc20add.TxID)
	if err != nil {
		return err
	}
	if hasTx.ID > 0 {
		return myerr.NewDbValidateError("txid已存在")
	}
	return t.Create()
}

//实收明细详情
func Erc20Detail(id int) (res response.Erc20List, err error) {
	ttx := model.TTxErc20{}
	ttx.ID = int64(id)
	res, err = ttx.GetErc20Detail()
	if err != nil {
		return
	}
	return
}

//根据商户订单编号绑定订单 解绑订单
func IsBindOrder(erc20 request.Erc20Edit) error {
	defer casbin.ClearEnforcer()
	//解绑
	if erc20.Status == 1 {
		err := model.BindErc20("", erc20.OrderId)
		if err != nil {
			return err
		}
		txErc20 := model.TTxErc20{}
		erc20Add := request.Erc20Edit{
			ID:      erc20.ID,
			OrderId: "",
			Status:  erc20.Status,
		}
		err1 := txErc20.BindOrder(erc20Add)
		if err1 != nil {
			return err1
		}
	} else {
		//查询订单明细
		ercdetail ,err:= model.GetErc20ByWhere("id = ?", erc20.ID)
		if err != nil {
			return err
		}
		//查询订单是否已经绑定
		hasTx ,err:= model.GetOrderByWhere("transaction_id = ?", ercdetail.SystemID);
		if err != nil {
			return err
		}
		if hasTx.ID > 0 {
			return myerr.NewDbValidateError("该订单已绑定")
		}
		err2 := model.BindErc20(ercdetail.SystemID, erc20.OrderId)
		if err2 != nil {
			return err2
		}
		txErc20 := model.TTxErc20{}
		erc20Add := request.Erc20Edit{
			ID:      erc20.ID,
			OrderId: erc20.OrderId,
			Status:  erc20.Status,
		}
		err1 := txErc20.BindOrder(erc20Add)
		if err1 != nil {
			return err1
		}
	}
	return nil
}
