package service
//
//import (
//	"j2pay-server/model"
//	"j2pay-server/model/request"
//	"j2pay-server/model/response"
//	"j2pay-server/pkg/casbin"
//	"j2pay-server/pkg/util"
//)
//
//// 实收明细记录列表
//func DetailedList(userId, status int, IdCode, address, txid, fromDate, toDate string, page, pageSize int) (res response.DetailedRecordPage, err error) {
//	detailed := model.DetailedRecord{}
//	if userId == 0 {
//		if status == 0 && fromDate == "" && toDate == "" && txid == "" && address == "" && IdCode == "" {
//			res, err = detailed.GetAll(page, pageSize, "pid = 0")
//		} else {
//			res, err = detailed.GetAll(page, pageSize, "status = ? or id_code like ? or txid like ? or address like ? or   or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", userId, status, IdCode, txid, address, fromDate, toDate)
//		}
//		if err != nil {
//			return
//		}
//	} else {
//		if status == 0 && fromDate == "" && toDate == "" && txid == "" && address == "" && IdCode == "" {
//			res, err = detailed.GetAll(page, pageSize, "user_id = ?", userId)
//		} else {
//			res, err = detailed.GetAll(page, pageSize, "user_id  = ? and status = ? or id_code like ? or txid like ? or address like ? or   or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", userId, status, IdCode, txid, address, fromDate, toDate)
//		}
//		if err != nil {
//			return
//		}
//	}
//	return
//}
//
////新增实收明细记录
//func DetailedAdd(detail request.DetailedAdd) error {
//	defer casbin.ClearEnforcer()
//	d := model.DetailedRecord{
//		IdCode:        util.RandString(20),
//		Amount:        detail.Amount,
//		TXID:          detail.TXID,
//		FromAddress:   detail.FromAddress,
//		ChargeAddress: detail.ChargeAddress,
//	}
//	return d.DetailedAdd()
//}
//
////实收明细详情
//func DetailsDetail(id int) (res response.DetailedList, err error) {
//	record := model.DetailedRecord{}
//	record.ID = int64(id)
//	res, err = record.GetDetail()
//	if err != nil {
//		return
//	}
//	return
//}
//
////绑定 解绑
//func IsBind(detail request.DetailedEdit)  error{
//	defer casbin.ClearEnforcer()
//	d := model.DetailedRecord{}
//	d.ID = int64(detail.ID)
//	//逻辑
//	return d.Binding(detail.ID,detail.OrderCode,detail.IsBind)
//}
